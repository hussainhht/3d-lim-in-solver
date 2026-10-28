// Vite bundles the actual repository maps as text, so no copied maps can drift.
import intro from "../../test/example00.txt?raw";
import parallel from "../../test/example01.txt?raw";
import direct from "../../test/example02.txt?raw";
import branching from "../../test/example05.txt?raw";
import colony from "../../3D/input.txt?raw";
import large from "../../test/example07.txt?raw";
export const examples = [
  {
    id: "colony",
    name: "The grand colony",
    tag: "SHOWCASE",
    description:
      "The native 3D demo map. Explore a branching colony and the routes selected by the real solver.",
    text: colony,
    file: "3D/input.txt",
  },
  {
    id: "intro",
    name: "First expedition",
    tag: "INTRODUCTORY",
    description:
      "Four ants, one corridor. A clear introduction to pipelined movement.",
    text: intro,
    file: "test/example00.txt",
  },
  {
    id: "parallel",
    name: "Parallel passages",
    tag: "MULTIPLE PATHS",
    description:
      "Ten ants and competing passages. Watch how the solver shares the colony across routes.",
    text: parallel,
    file: "test/example01.txt",
  },
  {
    id: "branching",
    name: "The labyrinth",
    tag: "BRANCHING",
    description:
      "Cross-connections and branching corridors challenge path selection.",
    text: branching,
    file: "test/example05.txt",
  },
  {
    id: "direct",
    name: "The shortcut",
    tag: "EDGE CASE",
    description:
      "A direct tunnel competes with a longer route. One ant may use a tunnel each turn.",
    text: direct,
    file: "test/example02.txt",
  },
  {
    id: "large",
    name: "A thousand travelers",
    tag: "LARGE COLONY",
    description:
      "The repository’s 1,000-ant example. A long playback on a compact graph.",
    text: large,
    file: "test/example07.txt",
  },
];
