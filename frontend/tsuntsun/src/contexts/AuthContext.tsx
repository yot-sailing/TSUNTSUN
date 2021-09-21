import axios from "axios";
import React, { useContext } from "react";

type auth = {
  isLoggedIn: () => boolean;
  idToken: () => string | null;
  accessToken: () => string | null;
  getloginHref: () => string;
  getToken: (code: string, string: string) => boolean;
};

const AuthContext = React.createContext<auth>({
  isLoggedIn: () => false,
  idToken: () => "",
  accessToken: () => "",
  getloginHref: () => "",
  getToken: () => false,
});

export const useAuth = () => {
  return useContext(AuthContext);
};

export const AuthProvider: React.FC = ({ children }) => {
  const isLoggedIn = () => localStorage.getItem("isLoggedIn") === "true";
  const idToken = () => localStorage.getItem("idToken");
  const accessToken = () => localStorage.getItem("accessToken");

  const getloginHref = () => {
    const state = Math.random().toString(32).substring(2);
    const nonce = Math.random().toString(32).substring(2);
    localStorage.setItem("state", state);
    localStorage.setItem("nonce", nonce);
    const href = `https://access.line.me/oauth2/v2.1/authorize?response_type=code&client_id=${process.env.REACT_APP_CHANNEL_ID}&redirect_uri=https://tsuntsun.herokuapp.com/after-login&state=${state}&scope=profile%20openid&nonce=${nonce}&bot_prompt=aggressive`;
    return href;
  };

  const getToken = (code: string, state: string): boolean => {
    // stateの確認
    const inputState = localStorage.getItem("state");
    console.log("login", inputState, state, code, isLoggedIn());

    if (inputState !== state) {
      localStorage.setItem("isLoggedIn", "false");
      return false;
    }
    // stateなど一時保存したものの削除
    localStorage.setItem("state", "");
    localStorage.setItem("nonce", "");

    // tokenの取得
    const data = {
      code: code,
      redirect_uri: "https://tsuntsun.herokuapp.com/after-login",
      client_id: process.env.REACT_APP_CHANNEL_ID,
      client_secret: process.env.REACT_APP_CHANNEL_SECRET,
    };
    axios
      .post("https://api.line.me/oauth2/v2.1/token", data, {
        headers: { "Content-Type": "application/x-www-form-urlencoded" },
      })
      .then((res) => {
        // tokenの保存
        localStorage.setItem("accessToken", res.data.access_token);
        localStorage.setItem("idToken", res.data.id_token);
        localStorage.setItem("isLoggedIn", "true");
        return true;
      })
      .catch((res) => console.log("catchres", res));
    return false;
  };

  const value = {
    isLoggedIn,
    idToken,
    accessToken,
    getloginHref,
    getToken,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};
