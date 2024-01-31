import { httpGet } from '@kfs/common/api/webServer';
import { useEffect, useState } from "react";
import { ScrollView, View } from "react-native";
import { Appbar, Divider, Text, useTheme } from "react-native-paper";
import Button from "../../components/Button";
import Drivers from './Drivers';

export default function () {
    const navigation = window.kfsNavigation;
    let [mediaTypes, setMediaTypes] = useState();
    let [locations, setLocations] = useState();
    const { colors } = useTheme();
    // console.log("mediaTypes", mediaTypes);
    // console.log("locations", locations);
    useEffect(() => {
        httpGet("/api/v1/listDCIMMediaType").then(setMediaTypes);
        httpGet("/api/v1/listDCIMLocation").then(setLocations);
    }, []);
    return (
        <View style={{ height: "100%" }}>
            <Appbar.Header mode="center-aligned">
                <Appbar.Content title="相册" />
            </Appbar.Header>

            <ScrollView style={{ flex: 1 }}>
                <Divider />
                <Text style={{ margin: 12, fontSize: 24, fontWeight: 400 }}>我的云盘</Text>
                <Drivers />

                <Divider />
                <Text style={{ margin: 12, fontSize: 24, fontWeight: 400 }}>地点</Text>
                <Button icon="map-outline" onPress={() => navigation.navigate("AlbumsLocation", { list: locations })} >
                    <Text style={{ color: colors.primary }}>地点</Text>
                    <Text style={{ color: colors.onSurfaceDisabled }}>{locations ? locations.length : "?"}</Text>
                </Button>

                <Divider />
                <Text style={{ margin: 12, fontSize: 24, fontWeight: 400 }}>人像</Text>

                <Button icon="face-recognition" disabled={true}>
                    人像
                </Button>

                <Divider />
                <Text style={{ margin: 12, fontSize: 24, fontWeight: 400 }}>物体</Text>
                <Button icon="flower-outline" disabled={true}>
                    物体
                </Button>

                <Divider />
                <Text style={{ margin: 12, fontSize: 24, fontWeight: 400 }}>媒体类型</Text>
                <Button icon="video-outline" onPress={() => navigation.navigate("AlbumsVideo", { list: mediaTypes.video })}>
                    <Text style={{ color: colors.primary }}>视频</Text>
                    <Text style={{ color: colors.onSurfaceDisabled }}>{mediaTypes ? mediaTypes.video.length : "?"}</Text>
                </Button>
                <Button icon="account-box-outline" onPress={() => navigation.navigate("AlbumsSelfie", { list: mediaTypes.selfie })}>
                    <Text style={{ color: colors.primary }}>自拍</Text>
                    <Text style={{ color: colors.onSurfaceDisabled }}>{mediaTypes ? mediaTypes.selfie.length : "?"}</Text>
                </Button>
                <Button icon="flower" disabled={true}>
                    实况
                </Button>
                <Button icon="panorama-horizontal-outline" disabled={true}>
                    全景
                </Button>
                <Button icon="timelapse" disabled={true}>
                    延时摄影
                </Button>
                <Button icon="motion-outline" disabled={true}>
                    慢动作
                </Button>
                <Button icon="movie-open-outline" disabled={true}>
                    电影效果
                </Button>
                <Button icon="cellphone-screenshot" disabled={true}>
                    截屏
                </Button>
                <Button icon="record-circle-outline" disabled={true}>
                    录屏
                </Button>
                <Button icon="motion-play-outline" disabled={true}>
                    动图
                </Button>

                <Divider />
                <Text style={{ margin: 12, fontSize: 24, fontWeight: 400 }}>更多项目</Text>
                <Button icon="pound" disabled={true}>
                    重复项目
                </Button>
            </ScrollView>
        </View>
    );
}
