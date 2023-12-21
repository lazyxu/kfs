import { httpGet } from '@kfs/common/api/webServer';
import useSysConfig from '@kfs/common/hox/sysConfig';
import { useMaterial3Theme } from '@pchmn/expo-material3-theme';
import { HoxRoot } from "hox";
import { useEffect, useState } from 'react';
import { View } from 'react-native';
import {
  BottomNavigation, MD3LightTheme as DefaultTheme, MD3DarkTheme,
  MD3LightTheme, PaperProvider,
  Text
} from 'react-native-paper';
import Toast from 'react-native-toast-message';
import "./global";
import Photos from './pages/Photos';
import SystemConfig from './pages/Setting/SystemConfig';

const MusicRoute = () => <Text>Music</Text>;

const AlbumsRoute = () => {
  let [drivers, setDrivers] = useState([]);
  console.log("drivers", drivers);
  useEffect(() => {
    httpGet("/api/v1/drivers").then(setDrivers);
  }, []);
  return (
    <View>
      {drivers.map(driver => (
        <Text key={driver.name}>{driver.name}</Text>
      ))}
    </View>
  );
};

const RecentsRoute = () => <Text>Recents</Text>;

const NotificationsRoute = () => {
  const { sysConfig, setSysConfig } = useSysConfig();
  console.log("sysConfig", sysConfig);
  return (
    <View>
      <Text>{JSON.stringify(sysConfig, undefined, 2)}</Text>
    </View>
  );
};

function App1() {
  const [index, setIndex] = useState(0);
  const [routes] = useState([
    { key: 'music', title: '照片', focusedIcon: 'heart', unfocusedIcon: 'heart-outline' },
    { key: 'albums', title: 'Albums', focusedIcon: 'album' },
    { key: 'recents', title: 'Recents', focusedIcon: 'history' },
    { key: 'notifications', title: 'Notifications', focusedIcon: 'bell', unfocusedIcon: 'bell-outline' },
  ]);

  const renderScene = BottomNavigation.SceneMap({
    music: Photos,
    albums: AlbumsRoute,
    recents: SystemConfig,
    notifications: NotificationsRoute,
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
