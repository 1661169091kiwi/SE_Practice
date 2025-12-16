# 接口新增与变更说明

## 认证
- 登录返回字段
  - `token`: `<jwt_token>`（HMAC-SHA256 签名，环境变量 `JWT_SECRET`）
  - `user`: 用户信息

## 采集员命名空间
- 获取比赛详情
  - `GET /api/collector/{sportType}/events/{eventId}`
  - 响应：
    ```json
    { "code":200, "data": { "id":1, "name":"...", "teamA":"...", "teamB":"...", "status":"in_progress", "currentData": { "scores": { "teamA":0, "teamB":0 }, "lineups": { "teamA":[], "teamB":[] } } }, "message":"success" }
    ```
- 提交数据
  - `POST /api/collector/{sportType}/events/{eventId}/data`
  - 请求：
    ```json
    {
      "eventId": "1",
      "sportType": "football",
      "timestamp": "2025-12-12T12:00:00Z",
      "collectorName": "张三",
      "status": "in_progress",
      "data": {
        "scores": { "teamA": 1, "teamB": 0 },
        "lineups": {
          "teamA": [ { "id":"20230001", "position":"FW", "number":"9", "isStarting": true } ],
          "teamB": []
        }
      }
    }
    ```
  - 响应：
    ```json
    { "code":200, "data": { "sportType":"football", "matchId":1, "updatedScore":true, "lineupsUpserted":1 }, "message":"success" }
    ```

## 比赛评分接口
- `GET/POST /api/matches/{match_id}/ratings`
- POST 请求体：
  ```json
  { "student_id":"20230001", "rater_id":"20231234", "score":9.5, "comment":"great" }
  ```

## 比赛阵容接口
- `GET/POST /api/matches/{match_id}/lineups`
- POST 请求体：
  ```json
  { "items": [ { "team_id":1, "student_id":"20230001", "position":"FW", "is_starting":true, "jersey_number":"9" } ] }
  ```

## 路由与响应统一
- 成功统一文案：`"message": "success"`
- 兼容赛事订阅路径：`POST /api/events/subscribe` 与 `POST /api/user/subscribe`
