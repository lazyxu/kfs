import { httpGet } from '@kfs/common/api/webServer';
import useSysConfig from '@kfs/common/hox/sysConfig';
import { HoxRoot } from "hox";
import { useEffect, useState } from 'react';
import { StyleSheet, Text, View } from 'react-native';
import { BottomNavigation, PaperProvider } from 'react-native-paper';
import Toast from 'react-native-toast-message';
import "./global";
import Photos from './pages/Photos';

const MusicRoute = () => <Text>Music</Text>;

const AlbumsRoute = () => {
  let [drivers, setDrivers] = useState([]);
  console.log("drivers", drivers);
  useEffect(() => {
    httpGet("/api/v1/drivers").then(setDrivers);
  }, []);
  return (
    <View style={styles.container}>
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
    <View style={styles.container}>
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
    recents: RecentsRoute,
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

export default function App() {
  return (
    <HoxRoot>
      <PaperProvider>
        <App1 />
        <Toast />
      </PaperProvider>
    </HoxRoot>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
    alignItems: 'center',
    justifyContent: 'center',
    overflow: "scroll"
  },
});
