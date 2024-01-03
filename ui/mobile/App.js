import { httpGet } from '@kfs/common/api/webServer';
import useSysConfig from '@kfs/common/hox/sysConfig';
import { useMaterial3Theme } from '@pchmn/expo-material3-theme';
import { NavigationContainer } from '@react-navigation/native';
import { createNativeStackNavigator } from '@react-navigation/native-stack';
import { HoxRoot } from "hox";
import * as React from 'react';
import { useEffect, useState } from 'react';
import { View } from 'react-native';
import {
  Appbar,
  BottomNavigation, MD3LightTheme as DefaultTheme, MD3DarkTheme,
  MD3LightTheme, PaperProvider,
  Text
} from 'react-native-paper';
import Toast from 'react-native-toast-message';
import "./global";
import Photos from './pages/Photos';
import SystemConfig from './pages/Setting/SystemConfig';
import ImageVideoViewer from './pages/Windows/ImageVideoViewer';
import Info from './pages/Windows/Info';

const Albums = () => {
  let [drivers, setDrivers] = useState([]);
  console.log("drivers", drivers);
  useEffect(() => {
    httpGet("/api/v1/drivers").then(setDrivers);
  }, []);
  return (
    <>
      <Appbar.Header mode="center-aligned">
        <Appbar.Content title="相册" />
      </Appbar.Header>
      <View style={{
        flex: 1,
        alignItems: 'center',
        justifyContent: 'center'
      }}>
        {drivers.map(driver => (
          <Text key={driver.name}>{driver.name}</Text>
        ))}
      </View>
    </>
  );
};

const Footprints = () => {
  return (
    <>
      <Appbar.Header mode="center-aligned">
        <Appbar.Content title="搜索" />
      </Appbar.Header>
      <View style={{
        flex: 1,
        alignItems: 'center',
        justifyContent: 'center'
      }}>
        <Text >TODO</Text>
      </View>
    </>
  );
};

function App1({ navigation }) {
  window.kfsNavigation = navigation;
  const [index, setIndex] = useState(0);
  const [routes] = useState([
    { key: 'photos', title: '照片', focusedIcon: 'image', unfocusedIcon: 'image-outline' },
    { key: 'albums', title: '相册', focusedIcon: 'image-multiple', unfocusedIcon: 'image-multiple-outline' },
    { key: 'footprints', title: '搜索', focusedIcon: 'image-search', unfocusedIcon: 'image-search-outline' },
    { key: 'me', title: '我', focusedIcon: 'account-settings', unfocusedIcon: 'account-settings-outline' },
  ]);

  const renderScene = BottomNavigation.SceneMap({
    photos: Photos,
    albums: Albums,
    footprints: Footprints,
    me: SystemConfig,
  });

  return (
    <BottomNavigation
      navigationState={{ index, routes }}
      onIndexChange={setIndex}
      renderScene={renderScene}
    />
  );
}

const theme = {
  ...DefaultTheme,
  // Specify custom property
  myOwnProperty: true,
  // Specify custom property in nested object
  colors: {
    ...DefaultTheme.colors,
    myOwnColor: '#BADA55',
  },
};

const Stack = createNativeStackNavigator();

function ThemeApp() {
  const { sysConfig } = useSysConfig();
  const { theme } = useMaterial3Theme();

  const paperTheme =
    sysConfig.theme === 'dark'
      ? { ...MD3DarkTheme, colors: theme.dark }
      : { ...MD3LightTheme, colors: theme.light };

  return (
    <PaperProvider theme={paperTheme}>
      <NavigationContainer>
        <Stack.Navigator initialRouteName="App1" screenOptions={{ headerShown: false, animation: 'slide_from_bottom' }}>
          <Stack.Screen name="App1" component={App1} />
          <Stack.Screen name="Viewer" component={ImageVideoViewer} />
          <Stack.Screen name="Info" component={Info} />
        </Stack.Navigator>
        <Toast />
      </NavigationContainer>
    </PaperProvider>
  );
}

function LoadingApp() {
  const { sysConfig } = useSysConfig();
  if (sysConfig) {
    return <ThemeApp />
  }
  return (
    <View style={{
      flex: 1,
      alignItems: 'center',
      justifyContent: 'center'
    }}>
      <Text >Loading...</Text>
    </View>
  );
}

export default function App() {
  return (
    <HoxRoot>
      <LoadingApp />
    </HoxRoot>
  );
}
