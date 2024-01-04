import { useCallback, useEffect, useRef, useState } from 'react';
import { Animated, Easing, Image, PanResponder, Platform, View } from 'react-native';
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

const DOUBLE_TAP_DELAY = 300; // milliseconds
const DOUBLE_TAP_RADIUS = 20;
function distance(x0, y0, x1, y1) {
    return Math.sqrt(Math.pow((x1 - x0), 2) + Math.pow((y1 - y0), 2));
}
function isDoubleTap(prev, cur) {
    return (cur.t - prev.t < DOUBLE_TAP_DELAY && distance(prev.x, prev.y, cur.x, cur.y) < DOUBLE_TAP_RADIUS);
}

export default function ({ navigation, route }) {
    const { hash, index, list } = route.params;
    const [hideHeaderFooter, setHideHeaderFooter] = useState(false);

    const [image, setImage] = useState();
    const [curIndex, setCurIndex, getCurIndex] = useGetState(index);
    const screenLayout = useRef();
    // console.log("ImageVideoViewer", curIndex, list);

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
                setImage({ origin, uri, x: 0, ...resize(origin, screenLayout.current), layout: screenLayout.current });
            }, err => {
                console.error(err);
                window.noteError("获取图片大小失败");
            });
        })()
    }, [curIndex]);

    const fadeAnim = useRef(new Animated.Value(1)).current;
    const [doubleTapGs, setDoubleTapGs] = useState();
    const isZoom = useRef(false);
    const onTap = () => setHideHeaderFooter(prev => !prev);
    useEffect(()=> {
        if (!doubleTapGs) {
            return;
        }
        console.log("useEffect.doubleTapGs", isZoom.current);
        if (!isZoom.current) {
            Animated.timing(fadeAnim, {
                toValue: 2,
                duration: 500,
                useNativeDriver: false,
            }).start();
        } else {
            Animated.timing(fadeAnim, {
                toValue: 1,
                duration: 500,
                useNativeDriver: false,
            }).start();
        }
        isZoom.current = !isZoom.current;
    }, [doubleTapGs]);
    const onDoubleTap = (gs) => {
        setDoubleTapGs({...gs});
        console.log("doubleTap", gs.x0, gs.y0);
    };
    const onNext = () => {
        if (getCurIndex() === list.length - 1) {
            setImage(img => ({ ...img, x: 0 }));
        } else {
            setCurIndex(i => i + 1);
        }
    }
    const onBack = () => {
        if (getCurIndex() === 0) {
            setImage(img => ({ ...img, x: 0 }));
        } else {
            setCurIndex(i => i - 1);
        }
    }

    const timer = useRef({});
    const prevent = useRef({});
    const prevTouchInfo = useRef({});
    const panResponder = useRef(PanResponder.create({
        onStartShouldSetPanResponder: () => true,
        onMoveShouldSetPanResponder: () => true,
        onPanResponderGrant: () => {
            // console.log('开始移动：');
        },
        onPanResponderMove: (evt, gs) => {
            // console.log('正在移动：X轴：' + gs.dx + '，Y轴：' + gs.dy);
            setImage(img => ({ ...img, x: gs.dx }));
        },
        onPanResponderRelease: (evt, gs) => {
            if (gs.dx > 50) {
                onBack();
            } else if (gs.dx < -50) {
                onNext();
            } else {
                let i = { x: gs.x0, y: gs.y0, t: Date.now() };
                if (isDoubleTap(prevTouchInfo.current, i)) {
                    clearTimeout(timer.current);
                    prevent.current = true;
                    onDoubleTap(gs);
                } else {
                    prevent.current = false;
                    timer.current = setTimeout(function () {
                        if (!prevent.current) {
                            onTap();
                        }
                        prevent.current = false;
                    }, DOUBLE_TAP_DELAY);
                    prevTouchInfo.current = i;
                }
            }
        }
    }));

    console.log("image", image);
    console.log("doubleTapGs", doubleTapGs);
    return (
        <>
            {!hideHeaderFooter && <Header navigation={navigation} hash={hash} index={curIndex} total={list.length}/>}
            <Surface style={{
                position: "absolute", left: 0, top: 0, right: 0, bottom: 0,
            }} onLayout={e => {
                const { layout } = e.nativeEvent;
                console.log("layout", layout)
                if (image) {
                    setImage(img => ({ ...img, ...resize(img.origin, layout), layout }));
                }
                screenLayout.current = layout;
            }} {...panResponder.current.panHandlers}>
                <View style={{
                    position: "absolute", left: image?.x,
                    width: "100%",
                    height: "100%",
                    overflow: 'hidden',
                    display: "flex",
                    flexDirection: "row",
                    justifyContent: "center",
                    alignItems: "center"
                }}>
                    {image ? <Animated.View
                        style={doubleTapGs && {
                            transform: [{
                                scale: fadeAnim,
                            }, {
                                translateX: fadeAnim.interpolate({
                                    inputRange: [1, 2],
                                    outputRange: [0, (image.layout.width/2-doubleTapGs.x0)/2],
                                }),
                            }, {
                                translateY: fadeAnim.interpolate({
                                    inputRange: [1, 2],
                                    outputRange: [0, (image.layout.height/2-doubleTapGs.y0)/2],
                                    easing: Easing.linear,
                                })
                            }]
                        }}>
                        <Image style={{
                            width: image.width,
                            height: image.height,
                        }} source={{ uri: image.uri }}
                        />
                    </Animated.View>
                        : <ActivityIndicator animating={true} size="large" />}
                </View>
            </Surface>
            {!hideHeaderFooter && <Footer navigation={navigation} hash={hash} uri={image?.uri}/>}
        </>
    );
}
