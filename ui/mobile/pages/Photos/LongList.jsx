import { memo, useState } from "react";
import { ScrollView, View } from "react-native";

export default memo(({ renderItem, itemHeightList, dataList }) => {
    const t0 = Date.now();
    const [curRect, setCurRect] = useState({ top: 0, bottom: 0 });
    const itemRects = [];
    for (let i = 0; i < itemHeightList.length; i++) {
        const h = itemHeightList[i];
        let top = 0;
        if (i !== 0) {
            top = itemRects[i - 1].bottom;
        }
        itemRects[i] = { top, bottom: top + h };
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
            inViewItems.push({ top: itemRects[i].top, key: i, elm: renderItem(dataList, i) });
        }
    }
    console.log("inView", itemRects.length, start, end);
    console.log("LongList.1", Date.now() - t0);
    return (
        <ScrollView style={{ height: "100%", width: "100%" }} scrollEventThrottle={0} onScroll={e => {
            setCurRect({
                top: e.nativeEvent.contentOffset.y,
                bottom: e.nativeEvent.contentOffset.y + e.nativeEvent.layoutMeasurement.height,
            });
            // console.log(itemRects, e.nativeEvent, e.nativeEvent.contentOffset.y)
        }} onLayout={e => {
            const { layout } = e.nativeEvent;
            if (layout.height) {
                setCurRect(prev => ({
                    top: prev.top,
                    bottom: prev.top + layout.height,
                }));
            }
        }}>
            <View style={{ height: itemRects.length > 0 ? itemRects[itemRects.length - 1].bottom : 0 }}>
                {inViewItems.map(item => <View key={item.key} style={{ position: "absolute", top: item.top }}>
                    {item.elm}
                </View>)}
            </View>
        </ScrollView >
    )
})
