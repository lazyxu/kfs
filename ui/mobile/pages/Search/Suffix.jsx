import { httpGet } from '@kfs/common/api/webServer';
import { useEffect, useState } from "react";
import { Appbar, Checkbox, List, Surface } from 'react-native-paper';
import { useCheckedSuffix } from '../../hox/checked';

export default function ({ navigation, route }) {
    // const { checked, setChecked } = route.params;
    const { checked, setChecked } = useCheckedSuffix();

    const [list, setList] = useState([]);

    useEffect(() => {
        httpGet("/api/v1/listDCIMSearchSuffix").then(setList);
    }, []);

    const ListItem = ({ item }) => {
        const key = item.suffix;
        const label = "." + item.suffix + " (" + item.count + ")";
        return <List.Item style={{ padding: 0, margin: 0 }} left={() =>
            <Checkbox.Item style={{ padding: 0, margin: 0, width: "100%" }} label={label} mode='ios'
                status={checked[key] ? 'checked' : 'unchecked'}
                onPress={() => {
                    console.log(route.params, checked, checked[key], checked[key] ? 0 : 1)
                    setChecked(m => {
                        return { ...m, [key]: m[key] ? 0 : 1 };
                    })
                }}
            />}
        />
    }

    return (
        <Surface style={{ height: "100%" }}>
            <Appbar.Header mode="center-aligned">
                <Appbar.BackAction onPress={() => navigation.pop()} />
                <Appbar.Content title="文件后缀" />
            </Appbar.Header>
            {list ? list.map(item => <ListItem item={item} />) :
                <View style={{
                    alignItems: 'center',
                    justifyContent: 'center'
                }}>
                    <Text >Loading...</Text>
                </View>}
        </Surface>
    );
};
