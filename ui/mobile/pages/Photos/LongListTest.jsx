import React, { useState } from "react";
import { View } from "react-native";
import { Text, useTheme } from "react-native-paper";
import LongList from "./LongList";
import Thumbnail from "./Thumbnail";

export default function ({ metadataTagList, elementsPerLine = 5, list, refresh }) {
    const { colors } = useTheme();
    const t0 = Date.now();
    const navigation = window.kfsNavigation;
    const [width, setWidth] = useState(0);
    // console.log("metadataTagList", metadataTagList, width)
    const itemHeightWidthList = [];
    let elementWidth = 0;
    if (width !== 0) {
        elementWidth = width / elementsPerLine;
        for (let i = 0; i < metadataTagList.length; i++) {
            const obj = metadataTagList[i];
            if (obj.hash) {
                itemHeightWidthList[i] = { height: elementWidth, width: elementWidth };
            } else {
                itemHeightWidthList[i] = { height: 20, width: width };
            }
        }
    }

    const renderItem = function (index, cacheItem) {
        const obj = metadataTagList[index];
        if (obj.hash) {
            return <Thumbnail key={obj.index} width={elementWidth} navigation={navigation} list={list} index={obj.index} inView={true} onLoaded={cacheItem} />;
        }
        cacheItem?.();
        return <View style={{
            height: "100%", width: "100%",
            display: "flex",
            flexDirection: "row",
            alignItems: "center",
            justifyContent: "space-between",
        }}>
            <Text style={{ color: colors.primary, fontWeight: 500, lineHeight: 20 }}>{obj.tag}</Text>
            <Text style={{ color: colors.onSurfaceDisabled, marginLeft: 4, lineHeight: 20 }}>{obj.end - obj.start + 1}</Text>
        </View>;
    }

    console.log("LongListTest.1", Date.now() - t0);
    return <View
        style={{ flex: 1 }}
        >
        {metadataTagList.length === 0 && <View style={{
            height: "100%", width: "100%",
            display: "flex",
            flexDirection: "row",
            alignItems: "center",
            justifyContent: 'center'
        }}>
            <Text style={{ color: colors.primary, fontWeight: 500, lineHeight: 20 }}>空空如也</Text>
        </View>}
        <LongList
            itemHeightWidthList={itemHeightWidthList}
            onWidth={setWidth}
            renderItem={renderItem}
            refresh={refresh}
        />
    </View>;
}
