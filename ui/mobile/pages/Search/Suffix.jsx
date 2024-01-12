import { httpGet } from '@kfs/common/api/webServer';
import * as React from 'react';
import { useEffect } from "react";
import { List, Text } from 'react-native-paper';

export default function () {
    const [expanded, setExpanded] = React.useState(false);
    const [list, setList] = React.useState([]);

    const handlePress = () => setExpanded(!expanded);

    useEffect(() => {
        httpGet("/api/v1/listDCIMSearchSuffix").then(setList);
    }, []);

    return (
        <List.Accordion
            title="文件后缀"
            left={props => <List.Icon {...props} icon="file-jpg-box" />}
            expanded={expanded}
            onPress={handlePress}>
            {list.map(item =>
                <List.Item
                    left={() => <Text style={{paddingLeft: 30}}>{"."+item.suffix}</Text>}
                    right={() => <Text>{item.count}</Text>}
                />
            )}
        </List.Accordion>
    );
};
