import { httpGet } from '@kfs/common/api/webServer';
import { useEffect, useState } from "react";
import { ScrollView, View } from 'react-native';
import { Appbar, Surface } from 'react-native-paper';
import { useCheckedSuffix } from '../../hox/checked';
import CheckItem from './CheckItem';

export default function ({ navigation, route }) {
    // const { checked, setChecked } = route.params;
    const [checked, setChecked] = useCheckedSuffix();

    const [list, setList] = useState([]);

    useEffect(() => {
        httpGet("/api/v1/listDCIMSearchSuffix").then(setList);
    }, []);

    const ListItem = ({ item }) => {
        const key = item.suffix;
        const label = "." + item.suffix + " (" + item.count + ")";
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
                <Appbar.Content title="文件后缀" />
            </Appbar.Header>
            <ScrollView style={{ flex: 1 }}>
                {list ? list.map((item, i) => <ListItem key={i} item={item} />) :
                    <View style={{
                        alignItems: 'center',
                        justifyContent: 'center'
                    }}>
                        <Text >Loading...</Text>
                    </View>}
            </ScrollView>
        </Surface>
    );
};
