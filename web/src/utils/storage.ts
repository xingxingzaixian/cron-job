const setStorageItem = (key: string, value: string, localType: string = 'localStorage') => {
  switch (localType) {
    case 'localStorage':
      localStorage.setItem(key, value);
      break;
    case 'sessionStorage':
      sessionStorage.setItem(key, value);
      break;
    default:
      localStorage.setItem(key, value);
      break;
  }
};

const getStorageItem = (key: string, localType: string = 'localStorage') => {
  switch (localType) {
    case 'localStorage':
      return localStorage.getItem(key);
    case 'sessionStorage':
      return sessionStorage.getItem(key);
    default:
      return localStorage.getItem(key);
  }
};

const removeStorageItem = (key: string, localType: string = 'localStorage') => {
  switch (localType) {
    case 'localStorage':
      localStorage.removeItem(key);
      break;
    case 'sessionStorage':
      sessionStorage.removeItem(key);
      break;
    default:
      localStorage.removeItem(key);
      break;
  }
};

const clearStorage = (localType: string = 'localStorage') => {
  switch (localType) {
    case 'localStorage':
      localStorage.clear();
      break;
    case 'sessionStorage':
      sessionStorage.clear();
      break;
    default:
      localStorage.clear();
      break;
  }
};

export { setStorageItem, getStorageItem, removeStorageItem, clearStorage };
