import { EventStreamContentType, fetchEventSource } from "react-native-fetch-event-source";
import { polyfill as polyfillFetch } from 'react-native-polyfill-globals/src/fetch';

polyfillFetch();

export { EventStreamContentType, fetchEventSource };

