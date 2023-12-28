import { listDCIMMetadataTime } from '@kfs/common/api/webServer/exif';
import { useEffect, useRef, useState } from "react";
import { View } from 'react-native';
import { Appbar, Surface, Text } from "react-native-paper";
import Thumbnail from './Thumbnail';

function calImageWith(gridWith) {
    return gridWith / 10;
}

export default function ({ navigation }) {
    const [metadataYearList, setMetadataYearList] = useState([]);
    const ref = useRef(null);
    const [width, setWidth] = useState(0);
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
            <View style={{
                flex: 1,
                overflow: "scroll",
            }} ref={ref}>
                {metadataYearList.map(metadataYear =>
                    <Surface key={metadataYear.year}>
                        <Text>{metadataYear.year === 1970 ? "未知时间" : metadataYear.year}</Text>
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
                    </Surface>
                )}
            </View>
        </View>
    );
}
