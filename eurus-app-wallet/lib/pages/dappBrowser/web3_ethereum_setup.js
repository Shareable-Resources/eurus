if (window.ethereum === undefined) {
  window.ethereum = {
    autoRefreshOnNetworkChange: false,
    isEurusWallet: true,

    request: async function (args) {
      const result = await window.flutter_inappwebview.callHandler(
        args.method,
        args.params
      );
      return new Promise((resolve, reject) => {
        if (
          typeof result === "object" &&
          result != null &&
          result.message != null
        ) {
          const error = Error(result.message);
          if (result.code != null) error.code = result.code;
          if (result.data != null) error.data = result.data;
          reject(error);
        } else {
          resolve(result);
        }
      });
    },

    send: function (method, params) {
      return this.request({ method: method, params: params });
    },

    on: function (method, handler) {
      switch (method) {
        case "chainChanged":
          this.chainChangedHandler = handler;
          break;
        case "accountsChanged":
          this.accountsChangedHandler = handler;
          break;
        case "networkChanged":
          this.networkChangedHandler = handler;
          break;
        case "message":
          this.messageHandler = handler;
          break;
        case "error":
          this.errorHandler = handler;
          break;
      }
    },
  };
  const { ethereum } = window;
}
