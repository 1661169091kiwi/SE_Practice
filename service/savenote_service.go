package service

import (
	"multimodal-notes/model"
	"multimodal-notes/repository"
	"strings"
)

// SaveNoteService 保存笔记的服务接口
type SaveNoteService interface {
	SaveNote(userID uint, originalText, summary, imageBase64 string, tags []string) (uint, error)
}

// saveNoteService 实现
type saveNoteService struct {
	noteRepo repository.NoteRepository
	tagRepo  repository.TagRepository
}

func NewSaveNoteService(noteRepo repository.NoteRepository, tagRepo repository.TagRepository) SaveNoteService {
	return &saveNoteService{
		noteRepo: noteRepo,
		tagRepo:  tagRepo,
	}
}

// extractTagsFromText 从文本中提取所有 # 开头的标签（ "a#b#c" → ["b", "c"]）
func extractTagsFromText(text string) []string {
	var tags []string
	if text == "" {
		return tags
	}

	// 按 #分割文本（分割后第一个元素是 #之前的内容，后续元素是标签候选）
	parts := strings.Split(text, "#")
	// 遍历分割后的元素，处理每个标签（去掉空格、空字符串）
	for _, part := range parts[1:] { // 从索引1开始，跳过 # 之前的内容
		// 去掉标签前后的空格（比如 "# 数学 " → "数学"）
		tag := strings.TrimSpace(part)
		// 过滤空标签（比如用户输入 "##"，分割后会得到空字符串）
		if tag != "" {
			tags = append(tags, tag)
		}
	}
	return tags
}

// SaveNote 保存笔记到数据库
func (s *saveNoteService) SaveNote(userID uint, originalText, summary, imageBase64 string, manualTags []string) (uint, error) {
	//这里的originalText是我让大模型生成标签并且直接以#XXX的形式直接写入原文开头，然后从中提取出来，后面的manualTags应该是前端设置一个专门的地方，用户在哪里手动输入tag
	extractedTags := extractTagsFromText(originalText)
	allTags := mergeAndDeduplicateTags(extractedTags, manualTags)

	tagList, err := s.tagRepo.GetOrCreateTags(allTags)
	if err != nil {
		return 0, err
	}
	note := &model.Note{
		Content:  originalText,
		Summary:  summary,
		ImageURL: imageBase64,
		UserID:   userID,  // 前端获取，把userID传进来
		Tags:     tagList, // 关联标签列表

	}
	return s.noteRepo.CreateNoteWithTags(note)
}
func mergeAndDeduplicateTags(list1, list2 []string) []string {
	tagMap := make(map[string]struct{}) // 用 map 去重
	for _, tag := range list1 {
		if tag != "" {
			tagMap[tag] = struct{}{}
		}
	}
	// list2用户手动传的标签
	for _, tag := range list2 {
		if tag != "" {
			tagMap[tag] = struct{}{}
		}
	}
	var result []string
	for tag := range tagMap {
		result = append(result, tag)
	}
	return result
}
