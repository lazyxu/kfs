import { httpGet } from '@kfs/common/api/webServer';
import useSysConfig from '@kfs/common/hox/sysConfig';
import { StatusBar } from 'expo-status-bar';
import { HoxRoot } from "hox";
import { useEffect, useState } from 'react';
import { StyleSheet, Text, View } from 'react-native';
import Toast from 'react-native-toast-message';
import "./global";
import Photos from './pages/Photos';

function App1() {
  const { sysConfig, setSysConfig } = useSysConfig();
  let [drivers, setDrivers] = useState([]);
  console.log("sysConfig", sysConfig);
  console.log("drivers", drivers);
  useEffect(() => {
    httpGet("/api/v1/drivers").then(setDrivers);
  }, []);
  return (
    <View style={styles.container}>
      <Text>App.js!</Text>
      <Text>{JSON.stringify(sysConfig, undefined, 2)}</Text>
      {drivers.map(driver => (
        <Text key={driver.name}>{driver.name}</Text>
      ))}
      <Photos />
      <StatusBar style="auto" />
      <Toast />
    </View>
  );
}

export default function App() {
  return (
    <HoxRoot>
      <App1 />
    </HoxRoot>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#fff',
    alignItems: 'center',
    justifyContent: 'center',
  },
});
