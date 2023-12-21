import { createNativeStackNavigator } from '@react-navigation/native-stack';
import * as React from 'react';
import Photos from './Photos';
import Viewer from './Viewer';

const Stack = createNativeStackNavigator();

export default function ({ navigation }) {
    return (
        <Stack.Navigator>
            <Stack.Screen name="Photos" component={Photos} />
            <Stack.Screen name="Viewer" component={Viewer} />
        </Stack.Navigator>
    );
}
