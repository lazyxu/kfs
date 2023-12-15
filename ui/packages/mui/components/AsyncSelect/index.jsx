import Autocomplete from '@mui/material/Autocomplete';
import CircularProgress from '@mui/material/CircularProgress';
import TextField from '@mui/material/TextField';
import { useEffect, useState } from 'react';

export default function AsyncSelect({ label, fetchOptions, onChange }) {
  const [open, setOpen] = useState(false);
  const [fetch, setFetch] = useState(false);
  const [options, setOptions] = useState([]);
  const loading = open && !fetch;

  useEffect(() => {
    let active = true;

    if (!loading) {
      return undefined;
    }

    (async () => {
      if (active) {
        setFetch(true);
        setOptions(await fetchOptions());
      }
    })();

    return () => {
      active = false;
    };
  }, [loading]);

  useEffect(() => {
    if (!open) {
      setFetch(false);
      setOptions([]);
    }
  }, [open]);

  return (
    <Autocomplete
      autowidth="true"
      open={open}
      onOpen={() => {
        setOpen(true);
      }}
      onClose={() => {
        setOpen(false);
      }}
      options={options}
      loading={loading}
      onChange={e => onChange(e.target.textContent)}
      renderInput={(params) => (
        <TextField
          {...params}
          label={label}
          InputProps={{
            ...params.InputProps,
            endAdornment: (
              <>
                {loading ? <CircularProgress color="inherit" size={20} /> : null}
                {params.InputProps.endAdornment}
              </>
            ),
          }}
        />
      )}
    />
  );
}
