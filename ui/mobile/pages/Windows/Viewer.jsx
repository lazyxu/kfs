import { default as React, useState } from 'react';
import { ActivityIndicator, Dimensions, View } from 'react-native';
import ImageViewer from 'react-native-image-zoom-viewer';
import { IconButton, Surface, useTheme } from "react-native-paper";
import ViewShot from 'react-native-view-shot';

const screenWidth = Dimensions.get("window").width;
const screenHeight = Dimensions.get("window").height;

export default function ({ navigation, route }) {
    const { hash, index, list } = route.params;
    console.log("Viewer", index, list);
    const [hideHeaderFooter, setHideHeaderFooter] = useState(false);
    const theme = useTheme();

    return (
        <>
            <ViewShot style={{
                padding: "0",
                width: "100%",
                height: "100%",
                display: "flex",
                justifyContent: "center",
                alignItems: "center"
            }}>
                <View style={{ position: "absolute", top: 0, bottom: 0, left: 0, right: 0 }}>
                    <ImageViewer
                        imageUrls={list} // 照片路径
                        enableImageZoom={true} // 是否开启手势缩放
                        // saveToLocalByLongPress={true} //是否开启长按保存
                        index={index} // 初始显示第几张
                        // failImageSource={} // 加载失败图片
                        backgroundColor={hideHeaderFooter ? null : theme.colors.surface}
                        renderIndicator={(currentIndex, allSize) =>
                            <></>
                        }
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
                                        // onPress={downloadImage}
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
                                    // onPress={shareImage}
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
                        // onSave={(url) => { this.savePhoto(url) }}
                    />
                </View>
            </ViewShot >
        </>
    );
}
