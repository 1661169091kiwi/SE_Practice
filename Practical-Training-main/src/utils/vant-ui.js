import { createApp } from 'vue';
import { Tabbar, TabbarItem,Icon,Search,Button } from 'vant';
import GameIcon from "@/assets/图标/赛事.svg"

const app = createApp();
app.use(Tabbar);
app.use(TabbarItem);
app.use(Icon);
app.use(Search);
app.use(Button);
app.use(GameIcon).provide("iconLoader",(name)=>{
  if(name==="GameIcon"){
    return GameIcon;
  }
  return null;
});
