import { listDCIMMetadataTime } from '@kfs/common/api/webServer/exif';
import { useCallback, useEffect, useRef, useState } from "react";
import { RefreshControl, ScrollView, View } from 'react-native';
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
        let yearList = [];
        let list;
        for (const m of l) {
            if (year !== m.year) {
                year = m.year;
                list = { year, list: [m.hash] }
                yearList.push(list);
            } else {
                list.list.push(m.hash);
            }
            // console.log(year, yearList, list, m)
        }
        setMetadataYearList(yearList);
    }
    const onRefresh = useCallback(() => {
      setRefreshing(true);
      refersh().then(() => {
        setRefreshing(false);
      });
    }, []);
    useEffect(() => {
        console.log("Photos useEffect");
        refersh();
        setWidth(calImageWith(ref.current.offsetWidth));
        return () => {
            console.log("Photos useEffect.unload");
        }
    }, []);
    return (
        <View style={{
            height: "100%",
            width: "100%",
            flexDirection: 'column',
        }}>
            <Appbar.Header mode="center-aligned">
                <Appbar.Content title="照片" />
                <Appbar.Action icon="calendar" onPress={() => { }} />
                <Appbar.Action icon="magnify" onPress={() => { }} />
            </Appbar.Header>
            <ScrollView
                showsVerticalScrollIndicator={true}
                style={{ flex: 1 }}
                stickyHeaderIndices={metadataYearList.map((_, i) => i * 2)}
                refreshControl={
                    <RefreshControl refreshing={refreshing} onRefresh={onRefresh} />
                }
                ref={ref}
            >
                {metadataYearList.map(metadataYear =>
                    [
                        <Surface><Text>{metadataYear.year === 1970 ? "未知时间" : metadataYear.year}</Text></Surface>,
                        <Surface style={{
                            display: "flex",
                            width: "100%",
                            flexDirection: 'row',
                            flexWrap: "wrap",
                            alignContent: "flex-start"
                        }}>
                            {metadataYear.list.map(hash =>
                                <Thumbnail key={hash} hash={hash} width={width} navigation={navigation} />
                            )}
                        </Surface>
                    ]
                )}
            </ScrollView>
        </View>
    );
}
