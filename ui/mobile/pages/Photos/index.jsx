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
    const [refreshing, setRefreshing] = useState(false);
    const refersh = async () => {
        const l = await listDCIMMetadataTime();
        let year = -1;
        let yearHashList = [];
        let hashList;
        for (const m of l) {
            if (year !== m.year) {
                year = m.year;
                hashList = [m.hash];
                yearHashList.push(year);
                yearHashList.push(hashList);
                // for (let i = 0; i < 1000; i++) {
                //     hashList.push(m.hash);
                // }
            } else {
                hashList.push(m.hash);
                // for (let i = 0; i < 1000; i++) {
                //     hashList.push(m.hash);
                // }
            }
            // console.log(year, yearList, list, m)
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
        setWidth(calImageWith(ref.current.offsetWidth));
        refersh();
        return () => {
            console.log("Photos useEffect.unload");
        }
    }, []);
    // console.log("render", width, navigation);
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
                showsVerticalScrollIndicator={true}
                style={{ flex: 1 }}
                stickyHeaderIndices={metadataYearList.filter((_, i) => i & 1 === 0).map((_, i) => i)}
                refreshControl={
                    <RefreshControl refreshing={refreshing} onRefresh={onRefresh} />
                }
                data={metadataYearList}
                extraData={width}
                renderItem={({ index, item }) => {
                    // console.log("render", index, index & 1 === 1, width, navigation, item);
                    return index & 1 === 1 ?
                        <Surface style={{
                            display: "flex",
                            width: "100%",
                            flexDirection: 'row',
                            flexWrap: "wrap",
                            alignContent: "flex-start"
                        }}>
                            {item.map((hash, i) =>
                                <Thumbnail key={i} hash={hash} width={width} navigation={navigation} />
                            )}
                        </Surface> :
                        <Surface><Text>{item === 1970 ? "未知时间" : item}</Text></Surface>
                }}
            />
        </View >
    );
}
