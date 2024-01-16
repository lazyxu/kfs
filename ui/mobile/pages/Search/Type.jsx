import { httpGet } from '@kfs/common/api/webServer';
import { useEffect, useState } from "react";
import { View } from 'react-native';
import { Appbar, Surface } from 'react-native-paper';
import { useCheckedType } from '../../hox/checked';
import CheckItem from './CheckItem';

export default function ({ navigation, route }) {
    // const { checked, setChecked } = route.params;
    const [checked, setChecked] = useCheckedType();

    const [list, setList] = useState([]);

    useEffect(() => {
        httpGet("/api/v1/listDCIMSearchType").then(setList);
    }, []);

    const Item = ({ item }) => {
        const key = item.type + "/" + item.subType;
        const label = item.type + "/" + item.subType + " (" + item.count + ")";
        return (
            <CheckItem label={label} status={checked[key] ? 'checked' : 'unchecked'}
                onPress={() => {
                    console.log(route.params, checked, checked[key], checked[key] ? 0 : 1)
                    setChecked(m => {
                        const ret = { ...m};
                        if (m[key]) {
                            delete ret[key];
                        } else {
                            ret[key] = 1;
                        }
                        return ret;
                    })
                }}
            />
        )
    }

    return (
        <Surface style={{ height: "100%" }}>
            <Appbar.Header mode="center-aligned">
                <Appbar.BackAction onPress={() => navigation.pop()} />
                <Appbar.Content title="文件类型" />
            </Appbar.Header>
            <View style={{ flex: 1, overflowY: "scroll" }}>
                {list ? list.map((item, i) => <Item key={i} item={item} />) :
                    <View style={{
                        alignItems: 'center',
                        justifyContent: 'center'
                    }}>
                        <Text >Loading...</Text>
                    </View>}
            </View>
        </Surface>
    );
};
