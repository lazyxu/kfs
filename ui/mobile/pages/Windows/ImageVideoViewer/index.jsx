import { useCallback, useEffect, useRef, useState } from 'react';
import { Image, PanResponder, Platform, View } from 'react-native';
import { CacheManager } from "react-native-expo-image-cache";
import {
    ActivityIndicator,
    Surface
} from 'react-native-paper';
import Footer from './Footer';
import Header from './Header';

function resize(origin, layout) {
    if (!layout) {
        return { width: 0, height: 0 };
    }
    let width = origin.width;
    let height = origin.height;
    if (width > layout.width) {
        const widthPixel = layout.width / width;
        width *= widthPixel;
        height *= widthPixel;
    }
    if (height > layout.height) {
        const HeightPixel = layout.height / height;
        width *= HeightPixel;
        height *= HeightPixel;
    }
    return { width, height };
}

const useGetState = (initiateState) => {
    const [state, setState] = useState(initiateState);
    const stateRef = useRef(state);
    stateRef.current = state;
    const getState = useCallback(() => stateRef.current, []);
    return [state, setState, getState];
};

export default function ({ navigation, route }) {
    const { hash, index, list } = route.params;
    const [hideHeaderFooter, setHideHeaderFooter] = useState(false);

    const [image, setImage] = useState();
    const [curIndex, setCurIndex, getCurIndex] = useGetState(index);
    const screenLayout = useRef();
    console.log("ImageVideoViewer", curIndex, list);

    useEffect(() => {
        (async () => {
            setImage();
            const origin = {};
            let uri = list[curIndex].url;
            origin.uri = uri;
            if (Platform.OS !== 'web') {
                uri = await CacheManager.get(uri).getPath();
                console.log("Cached", uri);
            }
            Image.getSize(uri, (width, height) => {
                origin.width = width;
                origin.height = height;
                setImage({ origin, uri, x: 0, ...resize(origin, screenLayout.current) });
            }, err => {
                console.error(err);
                window.noteError("获取图片大小失败");
            });
        })()
    }, [curIndex]);

    const panResponder = useRef(PanResponder.create({
        onStartShouldSetPanResponder: () => true,
        onMoveShouldSetPanResponder: () => true,
        onPanResponderGrant: () => {
            // console.log('开始移动：');
        },
        onPanResponderMove: (evt, gs) => {
            console.log('正在移动：X轴：' + gs.dx + '，Y轴：' + gs.dy);
            setImage(img => ({ ...img, x: gs.dx }));
        },
        onPanResponderRelease: (evt, gs) => {
            if (gs.dx > 50) {
                if (getCurIndex() === 0) {
                    setImage(img => ({ ...img, x: 0 }));
                } else {
                    setCurIndex(i => i - 1);
                }
            } else if (gs.dx < -50) {
                if (getCurIndex() === list.length - 1) {
                    setImage(img => ({ ...img, x: 0 }));
                } else {
                    setCurIndex(i => i + 1);
                }
            } else {
                setHideHeaderFooter(prev => !prev);
            }
        }
    }));

    console.log("image", image);
    return (
        <>
            {!hideHeaderFooter && <Header navigation={navigation} hash={hash} />}
            <Surface style={{
                position: "absolute", left: 0, top: 0, right: 0, bottom: 0,
            }} onLayout={e => {
                const { layout } = e.nativeEvent;
                if (image) {
                    setImage(img => ({ ...img, ...resize(img.origin, layout) }));
                }
                screenLayout.current = layout;
            }} {...panResponder.current.panHandlers}>
                <View style={{
                    position: "absolute", left: image?.x,
                    width: "100%",
                    height: "100%",
                    display: "flex",
                    flexDirection: "row",
                    justifyContent: "center",
                    alignItems: "center"
                }}>
                    {image ? <Image style={{
                        width: image.width,
                        height: image.height,
                    }} source={{ uri: image.uri }}
                    /> : <ActivityIndicator animating={true} size="large" />}
                </View>
            </Surface>
            {!hideHeaderFooter && <Footer navigation={navigation} hash={hash} />}
        </>
    );
}
