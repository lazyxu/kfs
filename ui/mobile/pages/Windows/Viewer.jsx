import { downloadURI } from "@kfs/common/api/web";
import { getSysConfig } from "@kfs/common/hox/sysConfig";
import * as Sharing from 'expo-sharing';
import { default as React, useEffect, useRef, useState } from 'react';
import { ActivityIndicator, Dimensions, Platform, View } from 'react-native';
import ImageViewer from 'react-native-image-zoom-viewer';
import { IconButton, Surface, useTheme } from "react-native-paper";
import ViewShot from 'react-native-view-shot';

const screenWidth = Dimensions.get("window").width;
const screenHeight = Dimensions.get("window").height;

export default function ({ navigation, route }) {
    const { hash } = route.params;
    console.log("Viewer", hash);
    const [downloadName, setDownloadName] = useState();
    const [metadata, setMetadata] = useState();
    const [openMetadata, setOpenMetadata] = useState(false);
    const [sameFiles, setSameFiles] = useState([]);
    const [openSameFiles, setOpenSameFiles] = useState(false);
    const [openAttribute, setOpenAttribute] = useState(false);
    const [hideHeaderFooter, setHideHeaderFooter] = useState(false);
    const theme = useTheme();

    const sysConfig = getSysConfig();
    const ref = useRef();
    const src = `${sysConfig.webServer}/api/v1/image?hash=${hash}`;

    const [blob, setBlob] = useState();
    const [url, setUrl] = useState();
    const controller = useRef(new AbortController());
    useEffect(() => {
        fetch(src, {
            method: 'get',
            signal: controller.current.signal,
        }).then(response => {
            return response.blob();
        }).then(b => {
            setBlob(b);
            const objectURL = URL.createObjectURL(b);
            console.log(b, objectURL);
            setUrl(objectURL);
        });
    }, []);

    const shareImage = async () => {
        try {
            console.log('shareImage', blob);
            if (Platform.OS === "web") {
                const b = new Blob([blob], { type: "image/png" });
                console.log(b);
                navigator.clipboard.write([
                    new ClipboardItem({
                        [b.type]: b
                    })
                ]);
                window.noteInfo(`拷贝图片成功：${hash}`);
                return
            }
            await Sharing.shareAsync({ url: url });
        } catch (e) {
            console.log(e);
        }
    };

    const downloadImage = async () => {
        try {
            console.log('downloadImage', blob);
            if (Platform.OS === "web") {
                const b = new Blob([blob], { type: "image/png" });
                downloadURI(url, `${hash}.png`);
                window.URL.revokeObjectURL(url);
                window.noteInfo(`下载图片成功：${hash}`);
                return
            }
            await Sharing.shareAsync({ url: url });
        } catch (e) {
            console.log(e);
        }
    };
    useEffect(() => {
        // getMetadata(hash).then(setMetadata);
        // listDriverFileByHash(hash).then(sf => {
        //     setSameFiles(sf);
        //     setDownloadName(sf[0].name);
        // });
    }, []);
    return (
        <>
            <ViewShot ref={ref} style={{
                padding: "0",
                width: "100%",
                height: "100%",
                display: "flex",
                justifyContent: "center",
                alignItems: "center"
            }}>
                {url && <View style={{ position: "absolute", top: 0, bottom: 0, left: 0, right: 0 }}><ImageViewer
                    imageUrls={[{ url }]} // 照片路径
                    enableImageZoom={true} // 是否开启手势缩放
                    saveToLocalByLongPress={true} //是否开启长按保存
                    index={0} // 初始显示第几张
                    // failImageSource={} // 加载失败图片
                    backgroundColor={hideHeaderFooter ? null : theme.colors.surface}
                    renderHeader={i => {
                        return !hideHeaderFooter && <Surface style={{
                            flexDirection: 'row',
                            alignItems: 'center',
                            justifyContent: "space-between",
                            paddingHorizontal: 4,
                            position: "fixed", left: 0, top: 0, right: 0,
                            zIndex: 1,
                        }}>
                            <IconButton
                                // size={size}
                                onPress={() => navigation.pop()}
                                // iconColor={actionIconColor}
                                icon="chevron-left"
                            // disabled={disabled}
                            // rippleColor={rippleColor}
                            />
                            <View style={{
                                flexDirection: 'row',
                                alignItems: 'center',
                            }}>
                                <IconButton
                                    // size={size}
                                    onPress={downloadImage}
                                    // iconColor={actionIconColor}
                                    icon="download-outline"
                                // disabled={disabled}
                                // rippleColor={rippleColor}
                                />
                                <IconButton
                                    // size={size}
                                    // onPress={downloadImage}
                                    // iconColor={actionIconColor}
                                    icon="dots-vertical"
                                // disabled={disabled}
                                // rippleColor={rippleColor}
                                />
                            </View>
                        </Surface>
                    }}
                    renderFooter={i => {
                        return !hideHeaderFooter && <Surface style={{
                            flexDirection: 'row',
                            alignItems: 'center',
                            justifyContent: "space-around",
                            paddingHorizontal: 4,
                            position: "fixed", left: 0, bottom: 0, right: 0,
                            zIndex: 1,
                        }}>
                            <IconButton
                                // size={size}
                                onPress={shareImage}
                                // iconColor={actionIconColor}
                                icon="export-variant"
                            // disabled={disabled}
                            // rippleColor={rippleColor}
                            />
                            <IconButton
                                // size={size}
                                onPress={() => navigation.navigate("Info", { hash })}
                                // iconColor={actionIconColor}
                                icon="information-outline"
                            // disabled={disabled}
                            // rippleColor={rippleColor}
                            />
                            <IconButton
                                // size={size}
                                // onPress={downloadImage}
                                // iconColor={actionIconColor}
                                icon="trash-can-outline"
                            // disabled={disabled}
                            // rippleColor={rippleColor}
                            />
                        </Surface>
                    }}
                    loadingRender={() => {
                        return <View style={{ marginTop: (screenHeight / 2) - 20 }}>
                            <ActivityIndicator animating={true} size="large" />
                        </View>
                    }}
                    enableSwipeDown={false}
                    menuContext={{ "saveToLocal": "保存图片", "cancel": "取消" }}
                    onChange={(index) => { }} // 图片切换时触发
                    onClick={() => { setHideHeaderFooter(prev => !prev) }}
                    onSave={(url) => { this.savePhoto(url) }}
                />
                </View>}
            </ViewShot >
        </>
    );
}
