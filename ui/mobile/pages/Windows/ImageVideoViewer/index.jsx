import { useEffect, useRef, useState } from 'react';
import { Image, PanResponder, Platform } from 'react-native';
import { CacheManager } from "react-native-expo-image-cache";
import {
    Surface,
    Text
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

export default function ({ navigation, route }) {
    const { hash, index, list } = route.params;
    const [hideHeaderFooter, setHideHeaderFooter] = useState(false);

    const [image, setImage] = useState();
    const [curIndex, setCurIndex] = useState(index);
    const screenLayout = useRef();
    console.log("ImageVideoViewer", curIndex, list);

    useEffect(() => {
        (async () => {
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
                setImage({ origin, uri, ...resize(origin, screenLayout.current) });
            }, err => {
                window.noteError(err.message);
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
            // console.log('正在移动：X轴：' + gs.dx + '，Y轴：' + gs.dy);
        },
        onPanResponderRelease: (evt, gs) => {
            if (gs.dx > 50) {
                setCurIndex(i => i === 0 ? i : i - 1);
            } else if (gs.dx < -50) {
                setCurIndex(i => i === list.length - 1 ? i : i + 1);
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
                display: "flex",
                flexDirection: "row",
                justifyContent: "center",
                alignItems: "center"
            }} onLayout={e => {
                const { layout } = e.nativeEvent;
                if (image) {
                    setImage(img => ({ ...img, ...resize(img.origin, layout) }));
                }
                screenLayout.current = layout;
            }} {...panResponder.current.panHandlers}>
                {image ? <Image style={{
                    width: image.width,
                    height: image.height,
                }} source={{ uri: image.uri }}
                /> : <Text>加载中...</Text>}
            </Surface>
            {!hideHeaderFooter && <Footer navigation={navigation} hash={hash} />}
        </>
    );
}
