import { httpGet } from '@kfs/common/api/webServer';
import { getSysConfig } from '@kfs/common/hox/sysConfig';
import { useEffect, useState } from "react";
import { Image, Pressable, View } from 'react-native';
import { Text } from "react-native-paper";

export default function () {
    let [drivers, setDrivers] = useState([]);
    console.log("drivers", drivers);
    useEffect(() => {
        httpGet("/api/v1/listDCIMDriver").then(setDrivers);
    }, []);
    return (
        <View style={{
            flexDirection: 'row',
            overflow: "scroll",
        }}>
            {drivers.map(driver => (
                <View key={driver.name} style={{
                    margin: 10,
                    marginRight: 0,
                }}>
                    <Pressable >
                        {driver.metadataList.length > 0 ? <Image style={{
                            height: 128,
                            width: 128,
                            borderRadius: 10,
                        }} source={{ uri: `${getSysConfig().webServer}/thumbnail?size=256&cutSquare=true&hash=${driver.metadataList[0].hash}` }} />
                            : <View style={{
                                height: 128,
                                width: 128,
                                borderRadius: 10,
                            }} />}
                        <Text>{driver.name}</Text>
                        <Text>{driver.metadataList.length}</Text>
                    </Pressable>
                </View>
            ))}
        </View>
    );
}
