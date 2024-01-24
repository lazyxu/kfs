import { useState } from "react";
import { View } from 'react-native';
import { Button } from "react-native-paper";
import ThumbnailListDay from './ThumbnailListDay';
import ThumbnailListMonth from './ThumbnailListMonth';
import ThumbnailListYear from './ThumbnailListYear';

const contentStyle = { justifyContent: 'flex-start' };

export default function ({ metadataList, listDCIMMetadataTime }) {
    const [value, setValue] = useState('年');
    const list = ["年", "月", "日"];
    return (
        <>
            {value === "年" && <ThumbnailListYear metadataList={metadataList} listDCIMMetadataTime={listDCIMMetadataTime} />}
            {value === "月" && <ThumbnailListMonth metadataList={metadataList} listDCIMMetadataTime={listDCIMMetadataTime} />}
            {value === "日" && <ThumbnailListDay metadataList={metadataList} listDCIMMetadataTime={listDCIMMetadataTime} />}
            <View style={{
                position: "absolute", bottom: 16, display: "flex", alignItems: "center", width: "100%", flexDirection: "row", justifyContent: "center",
                "div:first-child": {
                    borderRadius: 0,
                }
            }} >
                {list.map((name, i) => (
                    <Button onPress={() => setValue(name)}
                        mode={value === name ? "contained" : "elevated"}
                        style={i === 0 ? { borderTopRightRadius: 0, borderBottomRightRadius: 0 } :
                            i === list.length - 1 ? { borderTopLeftRadius: 0, borderBottomLeftRadius: 0 } :
                                { borderRadius: 0 }}
                    >
                        {name}
                    </Button>
                ))}
            </View >
        </>
    );
}
