import { httpGet } from '@kfs/common/api/webServer';
import { useEffect, useState } from "react";
import { Pressable } from "react-native";
import { Appbar, Surface, Text } from "react-native-paper";
import Drivers from './Drivers';

export default function () {
    const navigation = window.kfsNavigation;
    let [mediaTypes, setMediaTypes] = useState();
    let [locations, setLocations] = useState();
    console.log("mediaTypes", mediaTypes);
    console.log("locations", locations);
    useEffect(() => {
        httpGet("/api/v1/listDCIMMediaType").then(setMediaTypes);
        httpGet("/api/v1/listDCIMLocation").then(setLocations);
    }, []);
    return (
        <>
            <Appbar.Header mode="center-aligned">
                <Appbar.Content title="相册" />
            </Appbar.Header>

            <Surface><Text>我的云盘</Text></Surface>
            <Drivers />

            <Surface><Text>人物</Text></Surface>

            <Pressable onPress={() => navigation.navigate("AlbumsLocation", { list: locations })}>
                <Surface style={{
                    flexDirection: "row",
                    justifyContent: "space-between",
                }}>
                    <Text>地点</Text>
                    <Text>{locations ? locations.length : "?"}</Text>
                </Surface>
            </Pressable>

            <Surface><Text>媒体类型</Text></Surface>
            <Pressable disabled={!mediaTypes} onPress={() => navigation.navigate("AlbumsVideo", { list: mediaTypes.video })}>
                <Surface style={{
                    flexDirection: "row",
                    justifyContent: "space-between",
                }}>
                    <Text>视频</Text>
                    <Text>{mediaTypes ? mediaTypes.video.length : "?"}</Text>
                </Surface>
            </Pressable>
            <Pressable disabled={!mediaTypes} onPress={() => navigation.navigate("AlbumsSelfie", { list: mediaTypes.selfie })}>
                <Surface style={{
                    flexDirection: "row",
                    justifyContent: "space-between",
                }}>
                    <Text>自拍</Text>
                    <Text>{mediaTypes ? mediaTypes.selfie.length : "?"}</Text>
                </Surface>
            </Pressable>
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
