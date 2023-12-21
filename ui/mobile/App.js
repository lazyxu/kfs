import { httpGet } from '@kfs/common/api/webServer';
import useSysConfig from '@kfs/common/hox/sysConfig';
import { useMaterial3Theme } from '@pchmn/expo-material3-theme';
import { HoxRoot } from "hox";
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
        <Appbar.Content title="足迹" />
      </Appbar.Header>
      <View style={{
        flex: 1,
        alignItems: 'center',
        justifyContent: 'center'
      }}>
        TODO
      </View>
    </>
  );
};

function App1() {
  const [index, setIndex] = useState(0);
  const [routes] = useState([
    { key: 'photos', title: '照片', focusedIcon: 'image', unfocusedIcon: 'image-outline' },
    { key: 'albums', title: '相册', focusedIcon: 'image-multiple', unfocusedIcon: 'image-multiple-outline' },
    { key: 'footprints', title: '足迹', focusedIcon: 'image-marker', unfocusedIcon: 'image-marker-outline' },
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

function ThemeApp() {
  const { sysConfig } = useSysConfig();
  const { theme } = useMaterial3Theme();

  const paperTheme =
    sysConfig.theme === 'dark'
      ? { ...MD3DarkTheme, colors: theme.dark }
      : { ...MD3LightTheme, colors: theme.light };

  return (
    <PaperProvider theme={paperTheme}>
      <App1 />
      <Toast />
    </PaperProvider>
  );
}

export default function App() {
  return (
    <HoxRoot>
      <ThemeApp />
    </HoxRoot>
  );
}
