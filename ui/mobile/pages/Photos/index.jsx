import { listDCIMMetadataTime } from '@kfs/common/api/webServer/exif';
import { useState } from "react";
import { View } from 'react-native';
import { Appbar, SegmentedButtons, Surface } from "react-native-paper";
import ThumbnailList from './ThumbnailList';

export default function () {
    const [value, setValue] = useState('年');
    const [list, setList] = useState([]);
    return (
        <Surface style={{ height: "100%" }}>
            <Appbar.Header mode="center-aligned">
                <Appbar.Content title="照片" />
                <Appbar.Action icon="calendar" onPress={() => { }} />
                <Appbar.Action icon="magnify" onPress={() => { }} />
            </Appbar.Header>
            <ThumbnailList listDCIMMetadataTime={listDCIMMetadataTime} />
            <View style={{ position: "absolute", bottom: 16, display: "flex", alignItems: "center", width: "100%" }}>
                <SegmentedButtons
                    density="small"
                    value={value}
                    onValueChange={setValue}
                    buttons={[
                        {
                            value: "年",
                            label: "年",
                        },
                        {
                            value: "月",
                            label: "月",
                        },
                        {
                            value: "日",
                            label: "日",
                        },
                    ]}
                />
            </View>
        </Surface>
    );
}
