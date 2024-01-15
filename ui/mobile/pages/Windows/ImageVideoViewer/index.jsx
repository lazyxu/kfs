import { getSysConfig } from '@kfs/common/hox/sysConfig';
import { useCallback, useEffect, useRef, useState } from 'react';
import { Animated, Easing, Image, PanResponder, View } from 'react-native';
import {
    ActivityIndicator
} from 'react-native-paper';
import Footer from './Footer';
import Header from './Header';
import Video from './Video';

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
    const { list, index } = route.params;
    const [hideHeaderFooter, setHideHeaderFooter] = useState(false);

    const [image, setImage, getImage] = useGetState();
    const [curIndex, setCurIndex, getCurIndex] = useGetState(index);
    const screenLayout = useRef();
    // console.log("ImageVideoViewer", curIndex, list);

    useEffect(() => {
        (async () => {
            setImage();
            const origin = {};
            const hash = list[curIndex].hash;
            const type = list[curIndex].type;
            origin.height = list[curIndex].height;
            origin.width = list[curIndex].width;
            let uri = `${getSysConfig().webServer}/api/v1/image?hash=${hash}`;
            origin.uri = uri;
            console.log("origin", origin);
            // if (Platform.OS !== 'web') {
            //     uri = await CacheManager.get(uri).getPath();
            //     console.log("Cached", origin, uri);
            // }
            setImage({ origin, hash, type, uri, x: 0, ...resize(origin, screenLayout.current), layout: screenLayout.current });
        })()
    }, [curIndex]);

    const fadeAnim = useRef(new Animated.ValueXY(0, 0)).current;
    const [center, setCenter] = useState();
    const isZoom = useRef(false);
    const onTap = () => setHideHeaderFooter(prev => !prev);
    useEffect(() => {
        if (!center) {
            return;
        }
        console.log("useEffect.center", isZoom.current);
        console.log("layout", image.layout);
        console.log("center", center);
        Animated.timing(fadeAnim, {
            toValue: center,
            duration: 500,
            useNativeDriver: false,
            easing: Easing.linear,
        }).start();
    }, [center]);
    const onDoubleTap = (gs) => {
        console.log("doubleTap", gs.x0, gs.y0);
        if (!isZoom.current) {
            const layout = getImage().layout;
            setCenter({ x: (layout.width / 2 - gs.x0) / 2, y: (layout.height / 2 - gs.y0) / 2 });
        } else {
            Animated.timing(fadeAnim, {
                toValue: { x: 0, y: 0 },
                duration: 500,
                useNativeDriver: false,
                easing: Easing.linear,
            }).start();
        }
        isZoom.current = !isZoom.current;
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
    console.log("center", center);
    return (
        <>
            {!hideHeaderFooter && <Header navigation={navigation} hash={list[curIndex].hash} uri={image?.uri} index={curIndex} total={list.length} />}
            <View style={{
                position: "absolute", left: 0, top: 0, right: 0, bottom: 0,
            }} onLayout={e => {
                const { layout } = e.nativeEvent;
                console.log("layout", layout)
                if (getImage()) {
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
                        style={center && {
                            transform: [{
                                scale: fadeAnim.x.interpolate({
                                    inputRange: [0, Math.abs(center.x)],
                                    outputRange: [1, 2],
                                    easing: n => {
                                        // n: [0, -1]
                                        console.log(n, Math.pow(n, 2));
                                        // return [0, 1]
                                        return Math.abs(n);
                                    }
                                }),
                            }, {
                                translateX: fadeAnim.x,
                            }, {
                                translateY: fadeAnim.y
                            }]
                        }}>
                        {image.type === "image" && <Image style={{
                            width: image.width,
                            height: image.height,
                        }} source={{ uri: image.uri }}
                        />}
                        {image.type === "video" && <Video width={image.width} height={image.height} uri={image.uri} />}
                    </Animated.View>
                        : <ActivityIndicator animating={true} size="large" />}
                </View>
            </View>
            {!hideHeaderFooter && <Footer navigation={navigation} hash={list[curIndex].hash} uri={image?.uri} />}
        </>
    );
}
