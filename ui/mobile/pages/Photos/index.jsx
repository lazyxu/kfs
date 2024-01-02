import { listDCIMMetadataTime } from '@kfs/common/api/webServer/exif';
import { useCallback, useEffect, useRef, useState } from "react";
import { FlatList, RefreshControl, View } from 'react-native';
import { Appbar, Surface, Text } from "react-native-paper";
import Thumbnail from './Thumbnail';

function calImageWith(gridWith) {
    return gridWith / 10;
}

export default function () {
    const navigation = window.kfsNavigation;
    const [metadataYearList, setMetadataYearList] = useState([]);
    const ref = useRef(null);
    const [width, setWidth] = useState(0);
    const [initialNumToRender, setInitialNumToRender] = useState(0);
    const [refreshing, setRefreshing] = useState(false);
    const refersh = async () => {
        const l = await listDCIMMetadataTime();
        let year = -1;
        let yearHashList = [];
        let hashList;
        for (const m of l) {
            if (year !== m.year) {
                year = m.year;
                yearHashList.push(year);
                hashList = [m.hash];
                yearHashList.push(hashList);
                for (let i = 0; i < 100; i++) {
                    if (hashList.length == 10) {
                        hashList = [m.hash];
                        yearHashList.push(hashList);
                    } else {
                        hashList.push(m.hash);
                    }
                }
            } else {
                if (hashList.length == 10) {
                    hashList = [m.hash];
                    yearHashList.push(hashList);
                } else {
                    hashList.push(m.hash);
                }
                for (let i = 0; i < 100; i++) {
                    if (hashList.length == 10) {
                        hashList = [m.hash];
                        yearHashList.push(hashList);
                    } else {
                        hashList.push(m.hash);
                    }
                }
            }
            // console.log(yearHashList)
        }
        console.log("setMetadataYearList");
        setMetadataYearList(yearHashList);
    }
    const onRefresh = useCallback(() => {
        setRefreshing(true);
        refersh().then(() => {
            setRefreshing(false);
        });
    }, []);
    useEffect(() => {
        console.log("Photos useEffect");
        const w = ref.current.offsetWidth / 10;
        setWidth(w);
        setInitialNumToRender(Math.ceil(ref.current.offsetHeight / w));
        refersh();
        return () => {
            console.log("Photos useEffect.unload");
        }
    }, []);
    // console.log("render", width, navigation);
    let indices = [];
    for (let i = 0; i < metadataYearList.length; i++) {
        if (typeof metadataYearList[i] !== 'object') {
            indices.push(i);
        }
    }
    console.log(metadataYearList, indices, initialNumToRender)
    const renderItem = ({ index, item }) => {
        // console.log("render", index, index & 1 === 1, width, navigation, item);
        return typeof item === 'object' ?
            <View style={{
                display: "flex",
                width: "100%",
                flexDirection: 'row',
                flexWrap: "wrap",
                alignContent: "flex-start"
            }}>
                {item.map((hash, i) =>
                    <Thumbnail key={i} hash={hash} width={width} navigation={navigation} />
                )}
            </View> :
            <Surface><Text>{item === 1970 ? "未知时间" : item}</Text></Surface>
    };
    return (
        <View ref={ref} style={{
            height: "100%",
            width: "100%",
            flexDirection: 'column',
        }}>
            <Appbar.Header mode="center-aligned">
                <Appbar.Content title="照片" />
                <Appbar.Action icon="calendar" onPress={() => { }} />
                <Appbar.Action icon="magnify" onPress={() => { }} />
            </Appbar.Header>
            <FlatList
                // showsVerticalScrollIndicator={false}
                style={{ flex: 1 }}
                stickyHeaderIndices={indices}
                refreshControl={
                    <RefreshControl refreshing={refreshing} onRefresh={onRefresh} />
                }
                // renderScrollComponent
                initialNumToRender={initialNumToRender} // default 10
                maxToRenderPerBatch={1000} // default 10
                updateCellsBatchingPeriod={50} // default 50ms
                data={metadataYearList}
                extraData={width}
                renderItem={renderItem}
                // getItemLayout
            />
        </View >
    );
}
