import { useCallback, useEffect, useState } from "react";
import { FlatList, RefreshControl, View } from 'react-native';
import { Text } from "react-native-paper";
import Thumbnail from './Thumbnail';

export default function ({ metadataTagList, elementsPerLine = 5, list, refresh = () => { } }) {
    const navigation = window.kfsNavigation;
    const [width, setWidth] = useState(0);
    const [initialNumToRender, setInitialNumToRender] = useState(0);
    const [refreshing, setRefreshing] = useState(false);
    const onRefresh = useCallback(() => {
        setRefreshing(true);
        refresh().then(() => {
            setRefreshing(false);
        });
    }, []);
    useEffect(() => {
        refresh();
    }, []);
    // console.log("render", width, navigation);
    let indices = [];
    for (let i = 0; i < metadataTagList.length; i++) {
        if (typeof metadataTagList[i] !== 'object') {
            indices.push(i);
        }
    }
    console.log(metadataTagList, indices, initialNumToRender, list)
    console.log("width", width)
    console.log("initialNumToRender", initialNumToRender)
    const renderItem = ({ index, item }) => {
        // console.log("render", index, item, width);
        return typeof item === 'object' ?
            <View style={{
                display: "flex",
                width: "100%",
                flexDirection: 'row',
                flexWrap: "wrap",
                alignContent: "flex-start"
            }}>
                {item.map(({ hash, index }) =>
                    <Thumbnail key={index} width={width} navigation={navigation} list={list} index={index} />
                )}
            </View> :
            <View style={{ margin: 6 }}><Text style={{ fontSize: 16 }}>{item}</Text></View>
    };
    return (
        <FlatList
            onLayout={e => {
                const { layout } = e.nativeEvent;
                // console.log("onLayout", layout);
                if (layout.width) {
                    const w = layout.width / elementsPerLine;
                    setWidth(w);
                    setInitialNumToRender(Math.ceil(layout.height / w));
                }
            }}
            // showsVerticalScrollIndicator={false}
            style={{ flex: 1, width: "100%" }}
            // stickyHeaderIndices={indices}
            refreshControl={<RefreshControl refreshing={refreshing} onRefresh={onRefresh} />}
            // renderScrollComponent
            initialNumToRender={initialNumToRender} // default 10
            maxToRenderPerBatch={1000} // default 10
            updateCellsBatchingPeriod={50} // default 50ms
            data={metadataTagList}
            extraData={width}
            renderItem={renderItem}
        // getItemLayout
        />
    );
}
