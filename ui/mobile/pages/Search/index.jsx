import { useState } from "react";
import { View } from 'react-native';
import { Appbar, Button, SegmentedButtons } from 'react-native-paper';
import ThumbnailList from "../Photos/ThumbnailList";

const ItemStyle = { justifyContent: 'flex-start' };

export default function () {
    const navigation = window.kfsNavigation;
    const [value, setValue] = useState('filter');
    const [list, setList] = useState([]);
    return (
        <View style={{ height: "100%" }}>
            <Appbar.Header mode="center-aligned">
                <Appbar.Content title="搜索" />
            </Appbar.Header>
            <SegmentedButtons
                value={value}
                onValueChange={setValue}
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
            {value === "filter" && <View style={{ flex: 1, overflowY: "scroll" }}>
                <Button contentStyle={ItemStyle} icon="calendar-search" disabled={true}>
                    时间
                </Button>
                <Button contentStyle={ItemStyle} icon="map-search-outline" disabled={true}>
                    地点
                </Button>
                <Button contentStyle={ItemStyle} icon="file-image-outline" onPress={() => navigation.navigate("SearchType")}>
                    文件类型
                </Button>
                <Button contentStyle={ItemStyle} icon="text-recognition" disabled={true}>
                    文本识别
                </Button>

                <Button contentStyle={ItemStyle} icon="camera-outline" disabled={true}>
                    拍摄设备
                </Button>
                <Button contentStyle={ItemStyle} icon="numeric" disabled={true}>
                    文件大小
                </Button>
                <Button contentStyle={ItemStyle} icon="file-jpg-box" onPress={() => navigation.navigate("SearchSuffix")}>
                    文件后缀
                </Button>
            </View>}
            {value === "search" && <ThumbnailList metadataList={list} />}
        </View>
    );
}
