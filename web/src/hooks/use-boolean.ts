import { ref } from 'vue';

const useBoolean = (initValue = false) => {
  const bool = ref<boolean>(initValue);

  const setBool = (value: boolean) => {
    bool.value = value;
  };

  const setTrue = () => {
    setBool(true);
  };

  const setFalse = () => {
    setBool(false);
  };

  const toggle = () => {
    setBool(!bool.value);
  };

  return {
    bool,
    setBool,
    setTrue,
    setFalse,
    toggle
  };
};

export default useBoolean;
