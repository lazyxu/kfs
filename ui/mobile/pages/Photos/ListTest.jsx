/***
 Use this component inside your React Native Application.
 A scrollable list with different item type
 */
import React, { useCallback, useEffect, useRef, useState } from "react";
import { Dimensions, Text, View } from "react-native";
import Thumbnail from "./Thumbnail";
import { DataProvider, LayoutProvider, RecyclerListView } from "./recyclerlistview";

const ViewTypes = {
    FULL: 0,
    HALF_LEFT: 1,
    HALF_RIGHT: 2
};

let containerCount = 0;

class CellContainer extends React.Component {
    constructor(args) {
        super(args);
        this._containerId = containerCount++;
    }
    render() {
        return <View {...this.props}>{this.props.children}<Text>Cell Id: {this._containerId}</Text></View>;
    }
}

const _generateArray = (n) => {
    let arr = new Array(n);
    for (let i = 0; i < n; i++) {
        arr[i] = i;
    }
    return arr;
}

const provoder = new DataProvider((r1, r2) => {
    return r1 !== r2;
})

const useGetState = (initiateState) => {
    const [state, setState] = useState(initiateState);
    const stateRef = useRef(state);
    stateRef.current = state;
    const getState = useCallback(() => stateRef.current, []);
    return [state, setState, getState];
};

export default function ({ metadataTagList, elementsPerLine = 5, list, refresh }) {
    const navigation = window.kfsNavigation;
    let { width } = Dimensions.get("window");
    const [widthThumbnail, setWidthThumbnail] = useState(0);
    const [inViewIndexs, setInViewIndexs, getInViewIndexs] = useGetState([]);

    //Create the data provider and provide method which takes in two rows of data and return if those two are different or not.
    //THIS IS VERY IMPORTANT, FORGET PERFORMANCE IF THIS IS MESSED UP
    const [dataProvider, setDataProvider] = useState(provoder.cloneWithRows([""]));

    useEffect(() => {
        setDataProvider(provoder.cloneWithRows(metadataTagList.length === 0 ? ["无数据"] : metadataTagList));
    }, [metadataTagList]);

    // console.log(dataProvider, _generateArray(300), metadataTagList)
    //Create the layout provider
    //First method: Given an index return the type of item e.g ListItemType1, ListItemType2 in case you have variety of items in your list/grid
    //Second: Given a type and object set the exact height and width for that type on given object, if you're using non deterministic rendering provide close estimates
    //If you need data based check you can access your data provider here
    //You'll need data in most cases, we don't provide it by default to enable things like data virtualization in the future
    //NOTE: For complex lists LayoutProvider will also be complex it would then make sense to move it to a different file
    const _layoutProvider = new LayoutProvider(
        index => {
            return index;
        },
        (type, dim) => {
            const data = dataProvider.getDataForIndex(type);
            if (typeof data === 'object') {
                dim.width = width;
                dim.height = widthThumbnail;
            } else {
                dim.width = width;
                dim.height = 16 * 2;
            }
        }
    );


    //Given type and data return the view component
    const rowRenderer = (type, data, index) => {
        const inView = inViewIndexs.includes(index);
        console.log(type, data, index, inViewIndexs, inView);
        return typeof data === 'object' ?
            <View style={{
                display: "flex",
                width: "100%",
                flexDirection: 'row',
                flexWrap: "wrap",
                alignContent: "flex-start"
            }}>
                {data.map(({ hash, index }) =>
                    <Thumbnail key={index} width={widthThumbnail} navigation={navigation} list={list} index={index} inView={inView} />
                )}
            </View> :
            <View style={{ margin: 6 }}><Text style={{ fontSize: 16 }}>{data}</Text></View>
    }
    const onVisibleIndicesChanged = (all, now, notNow) => {
        setInViewIndexs(all);
        console.log("all,now, notNow", all, now, notNow);
    }
    console.log("inViewIndexs", inViewIndexs);

    return <View
        style={{ flex: 1 }}
        onLayout={e => {
            const { layout } = e.nativeEvent;
            // console.log("onLayout", layout);
            if (layout.width) {
                const w = layout.width / elementsPerLine;
                setWidthThumbnail(w);
                // setInitialNumToRender(Math.ceil(layout.height / w));
            }
        }}>
        <RecyclerListView
            layoutProvider={_layoutProvider}
            dataProvider={dataProvider}
            rowRenderer={rowRenderer}
            // extendedState={inViewIndexs}
            onVisibleIndicesChanged={onVisibleIndicesChanged}
        />
    </View>;
}

const styles = {
    container: {
        justifyContent: "space-around",
        alignItems: "center",
        flex: 1,
        backgroundColor: "#00a1f1"
    },
    containerGridLeft: {
        justifyContent: "space-around",
        alignItems: "center",
        flex: 1,
        backgroundColor: "#ffbb00"
    },
    containerGridRight: {
        justifyContent: "space-around",
        alignItems: "center",
        flex: 1,
        backgroundColor: "#7cbb00"
    }
};