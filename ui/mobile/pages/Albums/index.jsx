import { httpGet } from '@kfs/common/api/webServer';
import { useEffect, useState } from "react";
import { Appbar, Surface, Text } from "react-native-paper";
import Drivers from './Drivers';


export default function () {
    let [mediaTypes, setMediaTypes] = useState();
    console.log("mediaTypes", mediaTypes);
    useEffect(() => {
        httpGet("/api/v1/listDCIMMediaType").then(setMediaTypes);
    }, []);
    return (
      <>
        <Appbar.Header mode="center-aligned">
          <Appbar.Content title="相册" />
        </Appbar.Header>

        <Surface><Text>我的云盘</Text></Surface>
        <Drivers/>

        <Surface><Text>人物</Text></Surface>

        <Surface><Text>地点</Text></Surface>

        <Surface><Text>媒体类型</Text></Surface>
        <Surface><Text>视频 {mediaTypes?mediaTypes.video.length:0}</Text></Surface>
        <Surface><Text>自拍</Text></Surface>
        <Surface><Text>实况</Text></Surface>
        <Surface><Text>人像</Text></Surface>
        <Surface><Text>全景</Text></Surface>
        <Surface><Text>延时摄影</Text></Surface>
        <Surface><Text>慢动作</Text></Surface>
        <Surface><Text>电影效果</Text></Surface>
        <Surface><Text>截屏</Text></Surface>
        <Surface><Text>录屏</Text></Surface>
        <Surface><Text>动图</Text></Surface>

        <Surface><Text>重复项目</Text></Surface>
      </>
    );
}
