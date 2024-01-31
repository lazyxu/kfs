import { httpGet } from "@kfs/common/api/webServer";
import { useEffect, useState } from "react";
import { ScrollView, View } from 'react-native';
import { Appbar, SegmentedButtons, Text } from 'react-native-paper';
import Button from "../../components/Button";
import { useCheckedSuffix, useCheckedType } from "../../hox/checked";
import ThumbnailListAll from "../Photos/ThumbnailListAll";

export default function () {
    const navigation = window.kfsNavigation;
    const [value, setValue] = useState('filter');
    const [list, setList] = useState();
    const [checkedSuffix] = useCheckedSuffix();
    const [checkedType] = useCheckedType();

    const search = async () => {
        try {
            setList();
            const suffixList = Object.keys(checkedSuffix);
            const typeList = Object.keys(checkedType);
            console.log('api.searchImage', suffixList, typeList);
            const l = await httpGet("/api/v1/searchDCIM", {
                suffixList, typeList,
            });
            setList(l);
        } catch (e) {
            window.noteError("搜索失败：" + (typeof e.response?.data === 'string' ? e.response?.data : e.message));
        }
    }
    useEffect(() => {
        if (value === "search") {
            search();
        }
    }, [value]);
    return (
        <View style={{ height: "100%" }}>
            <Appbar.Header mode="center-aligned">
                <Appbar.Content title="搜索" />
            </Appbar.Header>
            <SegmentedButtons
                value={value}
                onValueChange={v => { setValue(v); setList(); }}
                buttons={[
                    {
                        value: 'filter',
                        icon: 'filter-outline',
                    },
                    {
                        value: 'search',
                        icon: 'magnify',
                    },
                ]}
            />
            {value === "filter" && <ScrollView style={{ flex: 1 }}>
                <Button icon="calendar-search" disabled={true}>
                    时间
                </Button>
                <Button icon="map-search-outline" disabled={true}>
                    地点
                </Button>
                <Button icon="file-image-outline" onPress={() => navigation.navigate("SearchType")}>
                    <Text>文件类型</Text>
                    <Text>{Object.keys(checkedType).length > 0 && `选中 ${Object.keys(checkedType).length} 项`}</Text>
                </Button>
                <Button icon="text-recognition" disabled={true}>
                    文本识别
                </Button>

                <Button icon="camera-outline" disabled={true}>
                    拍摄设备
                </Button>
                <Button icon="numeric" disabled={true}>
                    文件大小
                </Button>
                <Button icon="file-jpg-box" onPress={() => navigation.navigate("SearchSuffix")}>
                    <Text>文件后缀</Text>
                    <Text>{Object.keys(checkedSuffix).length > 0 && `选中 ${Object.keys(checkedSuffix).length} 项`}</Text>
                </Button>
            </ScrollView>}
            {value === "search" && (list ?
                <ThumbnailListAll metadataList={list} /> :
                <View style={{
                    alignItems: 'center',
                    justifyContent: 'center'
                }}>
                    <Text >Loading...</Text>
                </View>)}
        </View>
    );
}
