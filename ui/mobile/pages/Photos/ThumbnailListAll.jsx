import { useState } from "react";
import { View } from 'react-native';
import { SegmentedButtons } from "react-native-paper";
import ThumbnailListDay from './ThumbnailListDay';
import ThumbnailListMonth from './ThumbnailListMonth';
import ThumbnailListYear from './ThumbnailListYear';

export default function ({ metadataList, listDCIMMetadataTime }) {
    const [value, setValue] = useState('年');
    return (
        <>
            {value === "年" && <ThumbnailListYear metadataList={metadataList} listDCIMMetadataTime={listDCIMMetadataTime} />}
            {value === "月" && <ThumbnailListMonth metadataList={metadataList} listDCIMMetadataTime={listDCIMMetadataTime} />}
            {value === "日" && <ThumbnailListDay metadataList={metadataList} listDCIMMetadataTime={listDCIMMetadataTime} />}
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
        </>
    );
}
