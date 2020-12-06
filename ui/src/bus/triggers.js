import { basename } from 'utils/filepath';
import { addTrigger, setState, busState } from './bus';

addTrigger('chosen', (chosen) => {
  if (Object.values(chosen).filter((v) => v).length === 0) {
    setState({
      fileSize: null,
      boxChosen: {},
      mtimems: null,
      showDateInformations: false,
    });
  } else {
    try {
      setState({
        fileSize: Object.keys(chosen)
          .filter((k) => chosen[k])
          .map((k) => busState.files.find((f) => f.name === basename(k)).size)
          .reduce((a, b) => a + b, 0),
      });
    // eslint-disable-next-line no-empty
    } catch (e) {
    }
  }
});
