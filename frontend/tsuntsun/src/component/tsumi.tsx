import styled from "styled-components";
import { RiCalendar2Fill } from "react-icons/ri";
import Button from "./button";
import Tag from "./tag";
const mojs = require("@mojs/core");

export type TsumiObject = {
  author: string;
  category: string;
  createdAt?: string;
  deadline?: string;
  id: number;
  requiredTime: string;
  title: string;
  url: string;
  tags: TagObject[];
};

export type TagObject = {
  id: string;
  name: string;
};

const burst = new mojs.Burst({
  left: 0,
  top: 0,
  radius: { 0: 30 },
  angle: "rand(0, 360)",
  children: {
    shape: "line",
    stroke: "black",
    fill: "none",
    scale: 1,
    scaleX: { 1: 0 },
    easing: "cubic.out",
    duration: 1000,
  },
});

const bubbles = new mojs.Burst({
  left: 0,
  top: 0,
  radius: { 100: 200 },
  count: 40,
  timeline: { delay: 100 },
  children: {
    fill: [
      { "#DC93CF": "#E3D36B" },
      { "#91D3F7": "#9AE4CF" },
      { "#DC93CF": "#E3D36B" },
      { "#CF8EEF": "#CBEB98" },
      { "#F48BA2": "#CF8EEF" },
      { "#A7ECD0": "#9AE4CF" },
      { "#87E9C6": "#F48BA2" },
      { "#D58EB3": "#E0B6F5" },
      { "#F48BA2": "#F48BA2" },
      { "#91D3F7": "#A635D9" },
      { "#CF8EEF": "#CBEB98" },
      { "#87E9C6": "#F48BA2" },
    ],
    scale: { 1: 0, easing: "quad.in" },
    pathScale: [0.8, null],
    duration: [500, 700],
    easing: "quint.out",
    radius: { 0: "rand(6, 10)" },
    degreeShift: "rand(-50, 50)",
    delay: "rand(0, 250)",
  },
});

const Tsumi: React.FC<
  TsumiObject & { deleteFunc: (id: number) => void; isHist: boolean }
> = ({
  id,
  url,
  title,
  createdAt,
  requiredTime,
  deadline,
  tags,
  deleteFunc,
  isHist,
}) => {
  return (
    <Article key={id}>
      <ContentWrap>
        <Title>
          {url === "" ? title : <NoStyleLink href={url}>{title}</NoStyleLink>}
        </Title>
        <StatusBlack>
          <Status>
            {createdAt && subDate(createdAt) < 0 && (
              <>
                <NumBig>{-subDate(createdAt)}</NumBig>
                日経過
              </>
            )}
            {createdAt && subDate(createdAt) === 0 && <>本日</>}
          </Status>
          <Status>
            {requiredTime && (
              <>
                <NumBig>{requiredTime}</NumBig>
                分で読める
              </>
            )}
          </Status>
          <Status>
            {deadline && subDate(deadline) !== 0 && (
              <>
                {subDate(deadline) > 0 && "あと"}
                <NumBig>{Math.abs(subDate(deadline))}</NumBig>日
                {subDate(deadline) < 0 && "経過"}
              </>
            )}
            {deadline && subDate(deadline) === 0 && <p>本日</p>}
          </Status>
        </StatusBlack>
        <StatusBlack>
          {tags.map((t) => (
            <Tag key={t.id} name={t.name}></Tag>
          ))}
        </StatusBlack>
      </ContentWrap>
      <Button
        onClick={(e) => {
          deleteFunc(id);
          if (!isHist) {
            burst.tune({ x: e.pageX, y: e.pageY }).generate().replay();
            bubbles.tune({ x: e.pageX, y: e.pageY }).generate().replay();
          }
        }}
      >
        {isHist ? "戻す" : "読んだ"}
      </Button>
    </Article>
  );
};

const subDate: (date: string) => number = (date) => {
  const today = new Date();
  return Math.floor((Date.parse(date) - today.getTime()) / 86400000);
};

export default Tsumi;

const Title = styled.h3`
  text-align: left;
`;

const Article = styled.article`
  display: flex;
  align-items: center;
  max-width: 600px;
  padding: 20px 0;
  margin: 20px auto;
  border-bottom: 0.5px solid #b9b9b9;
`;

const NoStyleLink = styled.a`
  color: inherit;
  text-decoration: none;
`;

const StatusBlack = styled.div`
  display: flex;
`;

const Status = styled.div`
  display: inline-block;
  width: 30%;
  margin: 4px;
  text-align: left;
`;

const NumBig = styled.span`
  font-size: 24px;
  font-weight: 400;
  line-height: 28px;
`;

const ContentWrap = styled.div`
  flex: auto;
`;
