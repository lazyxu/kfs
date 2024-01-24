import { memo, useState } from "react";
import { ScrollView, View } from "react-native";

export default memo(({ renderItem, dataList, itemHeightWidthList, width }) => {
    const t0 = Date.now();
    const [curRect, setCurRect] = useState({ top: 0, bottom: 0 });
    const itemRects = [];
    let rect = { top: 0, left: 0, bottom: 0, right: 0 }
    for (let i = 0; i < itemHeightWidthList.length; i++) {
        const hw = itemHeightWidthList[i];
        let top, left;
        if (rect.right + hw.width > width) {
            top = rect.bottom;
            left = 0;
        } else {
            top = rect.top;
            left = rect.right;
        }
        itemRects[i] = { top, left, bottom: top + hw.height, right: left + hw.width };
        rect = itemRects[i];
    }
    let start = -1, end = -1;
    for (let i = 0; i < itemRects.length; i++) {
        if (start === -1 && curRect.top < itemRects[i].bottom) {
            start = i;
        }
        if (curRect.bottom < itemRects[i].top) {
            end = i - 1;
            break;
        }
    }
    if (end == -1) {
        end = itemRects.length - 1;
    }
    const inViewItems = [];
    if (itemRects.length !== 0) {
        for (let i = start; i <= end; i++) {
            inViewItems.push({ top: itemRects[i].top, left: itemRects[i].left, key: i, elm: renderItem(dataList, i) });
        }
    }
    console.log("inView", itemRects.length, start, end);
    console.log("LongList.1", Date.now() - t0);
    return (
        <ScrollView style={{ height: "100%", width: "100%" }} scrollEventThrottle={0} contentInsetAdjustmentBehavior="never"
            // contentContainerStyle={{ paddingRight: 14 }}
            onScroll={e => {
                setCurRect({
                    top: e.nativeEvent.contentOffset.y,
                    bottom: e.nativeEvent.contentOffset.y + e.nativeEvent.layoutMeasurement.height,
                });
                // console.log(itemRects, e.nativeEvent, e.nativeEvent.contentOffset.y)
            }} >
            <View style={{ height: rect.bottom }} onLayout={e => {
                const { layout } = e.nativeEvent;
                if (layout.height) {
                    setCurRect(prev => ({
                        top: prev.top,
                        bottom: prev.top + layout.height,
                    }));
                }
            }}>
                {inViewItems.map(item => <View key={item.key} style={{ position: "absolute", top: item.top, left: item.left }}>
                    {item.elm}
                </View>)}
            </View>
        </ScrollView >
    )
})
