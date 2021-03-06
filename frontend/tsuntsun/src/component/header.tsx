import styled from "styled-components";
import { useAuth } from "../contexts/AuthContext";
import { RiSettings3Fill, RiLogoutBoxRLine } from "react-icons/ri";
import AddButton from "./addButton";

const Header: React.FC<{
  name?: string;
}> = ({ name }) => {
  const auth = useAuth();
  const logout = async () => {
    await auth.logout();
    const isLoggedIn = await auth.isLoggedIn();
    if (!isLoggedIn) {
      window.location.href = "https://tsuntsun.herokuapp.com/login";
    }
  };
  return (
    <HeaderBase className="header">
      <Logo href="." className="brand">
        <img
          src={`${process.env.PUBLIC_URL}/logo.png`}
          alt=""
          style={{ height: "40px" }}
        ></img>
        <img
          src={`${process.env.PUBLIC_URL}/textonly.png`}
          alt=""
          style={{ height: "40px" }}
        ></img>
      </Logo>
      <Center>
        <AddButton />
        {name ? name + "さん、こんにちは" : ""}
        <RiSettings3Fill size={"2rem"} />
        <LogoutButton onClick={() => logout()}>
          <RiLogoutBoxRLine size={"2rem"} />
        </LogoutButton>
      </Center>
    </HeaderBase>
  );
};

export default Header;

const HeaderBase = styled.header`
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 65px;
  padding: 0 10%;
  background: #eaf5d3;
`;

const Center = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
`;

const Logo = styled.a`
  text-decoration: none;
`;

const LogoutButton = styled.button`
  background-color: transparent;
  border: none;
`;
