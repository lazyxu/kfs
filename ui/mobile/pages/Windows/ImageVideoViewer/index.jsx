import { useEffect, useRef, useState } from 'react';
import { Image, Platform, Pressable, View } from 'react-native';
import { CacheManager } from "react-native-expo-image-cache";
import {
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
    console.log("ImageVideoViewer", index, list);
    const [hideHeaderFooter, setHideHeaderFooter] = useState(false);

    const [image, setImage] = useState();
    const screenLayout = useRef();

    useEffect(() => {
        (async () => {
            const origin = {};
            let uri = list[index].url;
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
    }, []);
    console.log("image", image);
    return (
        <>
            {!hideHeaderFooter && <Header navigation={navigation} hash={hash} />}
            <Pressable onPress={() => { setHideHeaderFooter(prev => !prev) }}>
                <View style={{
                    position: "fixed",
                    width: "100%",
                    height: "100%",
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
                }}>
                    {image ? <Image style={{
                        width: image.width,
                        height: image.height,
                    }} source={{ uri: image.uri }}
                    /> : <Text>加载中...</Text>}
                </View>
            </Pressable>
            {!hideHeaderFooter && <Footer navigation={navigation} hash={hash} />}
        </>
    );
}
