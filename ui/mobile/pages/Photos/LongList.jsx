import { memo, useCallback, useEffect, useRef, useState } from "react";
import { ScrollView, View } from "react-native";

const useGetState = (initiateState) => {
    const [state, setState] = useState(initiateState);
    const stateRef = useRef(state);
    stateRef.current = state;
    const getState = useCallback(() => stateRef.current, []);
    return [state, setState, getState];
};

export default memo(({ renderItem, itemHeightWidthList, onWidth }) => {
    const t0 = Date.now();
    const [curRect, setCurRect, getCurRect] = useGetState({ top: 0, bottom: 0, width: 0 });
    const [itemRects, setItemRects, getItemRects] = useGetState([]);
    const cacheIndexes = useRef([]);
    useEffect(() => {
        const _itemRects = [];
        let rect = { top: 0, left: 0, bottom: 0, right: 0 };
        for (let i = 0; i < itemHeightWidthList.length; i++) {
            const hw = itemHeightWidthList[i];
            let top, left;
            if (rect.right + hw.width > getCurRect().width) {
                top = rect.bottom;
                left = 0;
            } else {
                top = rect.top;
                left = rect.right;
            }
            // console.log("_itemRects", i, { top, left, bottom: top + hw.height, right: left + hw.width });
            _itemRects[i] = { top, left, bottom: top + hw.height, right: left + hw.width };
            rect = _itemRects[i];
        }
        // console.log("_itemRects", _itemRects, itemHeightWidthList);
        setItemRects(_itemRects);
    }, [itemHeightWidthList]);
    const _curRect = getCurRect();
    const _itemRects = getItemRects();
    // console.log("_curRect", _curRect, _itemRects);
    let start = -1, end = -1;
    for (let i = 0; i < _itemRects.length; i++) {
        if (start === -1 && _curRect.top < _itemRects[i].bottom) {
            start = i;
        }
        if (_curRect.bottom < _itemRects[i].top) {
            end = i - 1;
            break;
        }
    }
    if (end == -1) {
        end = _itemRects.length - 1;
    }
    const inViewItems = {};
    if (_itemRects.length !== 0) {
        // console.log("cacheIndexes", cacheIndexes.current);
        for (const i of cacheIndexes.current) {
            inViewItems[i] = {
                top: _itemRects[i].top, left: _itemRects[i].left, right: _itemRects[i].right, bottom: _itemRects[i].bottom,
                key: i, elm: renderItem(i),
            };
        }
        for (let i = start; i <= end; i++) {
            if (inViewItems.hasOwnProperty(i)) {
                continue;
            }
            inViewItems[i] = {
                top: _itemRects[i].top, left: _itemRects[i].left, right: _itemRects[i].right, bottom: _itemRects[i].bottom,
                key: i, elm: renderItem(i, () => {
                    // console.log("cacheIndex", i);
                    cacheIndexes.current.push(i);
                })
            };
        }
    }
    console.log("inView", _itemRects.length, _curRect, start, end);
    console.log("cacheIndexes", cacheIndexes.current.length);
    console.log("inViewItems", inViewItems.length);
    console.log("LongList.1", Date.now() - t0);
    return (
        <ScrollView scrollEventThrottle={0} contentInsetAdjustmentBehavior="never"
            // contentContainerStyle={{ paddingRight: 14 }}
            onScroll={e => {
                setCurRect({
                    top: e.nativeEvent.contentOffset.y,
                    bottom: e.nativeEvent.contentOffset.y + e.nativeEvent.layoutMeasurement.height,
                });
                // console.log(itemRects, e.nativeEvent, e.nativeEvent.contentOffset.y)
            }}
            // onLayout={e => {
            //     const { layout } = e.nativeEvent;
            //     if (layout.height) {
            //         setCurRect(prev => ({
            //             top: prev.top,
            //             bottom: prev.top + layout.height,
            //             width: layout.width,
            //         }));
            //         onWidth(layout.width);
            //     }
            // }}
            onContentSizeChange={(w, h) => {
                console.log("onContentSizeChange", w, h)
                setCurRect(prev => ({
                    top: prev.top,
                    bottom: prev.top + h,
                    width: w,
                }));
                onWidth(w);
            }}
            >
            <View style={{ height: itemRects.length > 0 ? itemRects[itemRects.length - 1].bottom : 0 }}>
                {Object.values(inViewItems).map(item =>
                    <View key={item.key} style={{
                        position: "absolute", top: item.top, left: item.left,
                        width: item.right - item.left, height: item.bottom - item.top
                    }}
                    >
                        {item.elm}
                    </View>)}
            </View>
        </ScrollView >
    )
})
