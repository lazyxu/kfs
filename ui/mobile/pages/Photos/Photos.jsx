import { listDCIMMetadataTime } from '@kfs/common/api/webServer/exif';
import { getSysConfig } from "@kfs/common/hox/sysConfig";
import { useEffect, useRef, useState } from "react";
import { Image, TouchableHighlight, View } from 'react-native';
import { Appbar, Text } from "react-native-paper";
import Viewer from './Viewer';

function calImageWith(gridWith) {
    return gridWith / 10;
}

export default function ({ navigation }) {
    const [metadataYearList, setMetadataYearList] = useState([]);
    const ref = useRef(null);
    const [width, setWidth] = useState(0);
    const [openViewer, setOpenViewer] = useState();
    const refersh = () => {
        listDCIMMetadataTime().then(l => {
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
        });
    }
    useEffect(() => {
        refersh();
        setWidth(calImageWith(ref.current.offsetWidth));
    }, []);
    return (
        <>
            <Appbar.Header mode="center-aligned">
                <Appbar.Content title="照片" />
                <Appbar.Action icon="calendar" onPress={() => { }} />
                <Appbar.Action icon="magnify" onPress={() => { }} />
            </Appbar.Header>
            <View style={{
                width: "100%",
                height: "100%",
                overflow: "scroll",
            }} ref={ref}>
                {metadataYearList.map(metadataYear =>
                    <View key={metadataYear.year}>
                        <Text>{metadataYear.year === 1970 ? "未知时间" : metadataYear.year}</Text>
                        <View key={metadataYear.year} style={{
                            display: "flex",
                            backgroundColor: '#fff',
                            width: "100%",
                            flexDirection: 'row',
                            flexWrap: "wrap",
                            alignContent: "flex-start"
                        }}>
                            {metadataYear.list.map(hash =>
                                <TouchableHighlight key={hash} onPress={() => { navigation.navigate("Viewer", { hash }); }}  >
                                    <Image style={{
                                        width: width,
                                        height: width,
                                    }} source={{ uri: `${getSysConfig().webServer}/thumbnail?size=256&cutSquare=true&hash=${hash}` }} />
                                </TouchableHighlight>
                            )}
                        </View>
                    </View>
                )}
            </View>
            {openViewer && <Viewer hash={openViewer} onClose={() => { setOpenViewer(); }} />}
        </>
    );
}
