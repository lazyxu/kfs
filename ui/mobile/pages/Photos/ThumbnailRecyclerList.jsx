import { useCallback, useState } from "react";
import { View } from 'react-native';
import { Text } from "react-native-paper";
import { RecyclerListView } from "recyclerlistview/web";
import Thumbnail from './Thumbnail';

export default function ({ metadataTagList, elementsPerLine = 5, list, refresh }) {
    const navigation = window.kfsNavigation;
    const [width, setWidth] = useState(0);
    const [initialNumToRender, setInitialNumToRender] = useState(0);
    const [refreshing, setRefreshing] = useState(false);
    const onRefresh = useCallback(() => {
        setRefreshing(true);
        refresh().then(() => {
            setRefreshing(false);
        });
    }, []);
    // console.log("render", width, navigation);
    let indices = [];
    for (let i = 0; i < metadataTagList.length; i++) {
        if (typeof metadataTagList[i] !== 'object') {
            indices.push(i);
        }
    }
    console.log(metadataTagList, indices, initialNumToRender, list)
    console.log("width", width)
    console.log("initialNumToRender", initialNumToRender)
    const renderItem = ({ index, item }) => {
        // console.log("render", index, item, width);
        return typeof item === 'object' ?
            <View style={{
                display: "flex",
                width: "100%",
                flexDirection: 'row',
                flexWrap: "wrap",
                alignContent: "flex-start"
            }}>
                {item.map(({ hash, index }) =>
                    <Thumbnail key={index} width={width} navigation={navigation} list={list} index={index} />
                )}
            </View> :
            <View style={{ margin: 6 }}><Text style={{ fontSize: 16 }}>{item}</Text></View>
    };
    return (
        <RecyclerListView
            rowRenderer={this.rowRenderer}
            dataProvider={queue}
            layoutProvider={this.layoutProvider}
            onScroll={this.checkRefetch}
            // renderFooter={this.renderFooter}
            // scrollViewProps={{
            //     refreshControl: (
            //         <RefreshControl
            //             refreshing={loading}
            //             onRefresh={async () => {
            //                 this.setState({ loading: true });
            //                 analytics.logEvent('Event_Stagg_pull_to_refresh');
            //                 await refetchQueue();
            //                 this.setState({ loading: false });
            //             }}
            //         />
            //     )
            // }}
        />
    );
}
