import { getSysConfig } from "@kfs/common/hox/sysConfig";
import { useCallback, useEffect, useState } from "react";
import { FlatList, RefreshControl, View } from 'react-native';
import { Surface, Text } from "react-native-paper";
import Thumbnail from './Thumbnail';

export default function ({ listDCIMMetadataTime, metadataList }) {
    const navigation = window.kfsNavigation;
    const [metadataYearList, setMetadataYearList] = useState([]);
    const sysConfig = getSysConfig();
    const [width, setWidth] = useState(0);
    const [initialNumToRender, setInitialNumToRender] = useState(0);
    const [refreshing, setRefreshing] = useState(false);
    const [list, setList] = useState([]);
    const refersh = async () => {
        let l;
        if (metadataList) {
            l = metadataList;
        }
        if (listDCIMMetadataTime) {
            l = await listDCIMMetadataTime();
        }
        let year = -1;
        let yearHashList = [];
        let hashList;
        const allHashList = [];
        l = l.slice(0, 100);
        for (let index = 0; index < l.length; index++) {
            // console.log(index);
            const m = l[index];
            allHashList.push({
                url: `${sysConfig.webServer}/api/v1/image?hash=${m.hash}`,
                hash: m.hash,
                type: m.fileType.type,
                height: m.heightWidth.height,
                width: m.heightWidth.width,
            });
            if (year !== m.year) {
                year = m.year;
                yearHashList.push(year);
                hashList = [{ index, hash: m.hash }];
                yearHashList.push(hashList);
                // for (let i = 0; i < 100; i++) {
                //     allHashList.push({url: `${sysConfig.webServer}/api/v1/image?hash=${m.hash}`});
                //     if (hashList.length == 10) {
                //         hashList = [{ index, hash: m.hash }];
                //         yearHashList.push(hashList);
                //     } else {
                //         hashList.push({ index, hash: m.hash });
                //     }
                // }
            } else {
                if (hashList.length == 10) {
                    hashList = [{ index, hash: m.hash }];
                    yearHashList.push(hashList);
                } else {
                    hashList.push({ index, hash: m.hash });
                }
                // for (let i = 0; i < 100; i++) {
                //     allHashList.push({url: `${sysConfig.webServer}/api/v1/image?hash=${m.hash}`});
                //     if (hashList.length == 10) {
                //         hashList = [{ index, hash: m.hash }];
                //         yearHashList.push(hashList);
                //     } else {
                //         hashList.push({ index, hash: m.hash });
                //     }
                // }
            }
            // console.log(yearHashList)
        }
        console.log("setMetadataYearList");
        setList(allHashList);
        setMetadataYearList(yearHashList);
    }
    const onRefresh = useCallback(() => {
        setRefreshing(true);
        refersh().then(() => {
            setRefreshing(false);
        });
    }, []);
    useEffect(() => {
        refersh();
    }, []);
    // console.log("render", width, navigation);
    let indices = [];
    for (let i = 0; i < metadataYearList.length; i++) {
        if (typeof metadataYearList[i] !== 'object') {
            indices.push(i);
        }
    }
    console.log(metadataYearList, indices, initialNumToRender, list)
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
                    <Thumbnail key={index} hash={hash} width={width} navigation={navigation} list={list} index={index} />
                )}
            </View> :
            <Surface><Text>{item === 1970 ? "未知时间" : item}</Text></Surface>
    };
    return (
        <FlatList
            onLayout={e => {
                const { layout } = e.nativeEvent;
                // console.log("onLayout", layout);
                if (layout.width) {
                    const w = layout.width / 10;
                    setWidth(w);
                    setInitialNumToRender(Math.ceil(layout.height / w));
                }
            }}
            // showsVerticalScrollIndicator={false}
            style={{ flex: 1, width: "100%" }}
            // stickyHeaderIndices={indices}
            refreshControl={
                listDCIMMetadataTime ? <RefreshControl refreshing={refreshing} onRefresh={onRefresh} /> : undefined
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
    );
}
