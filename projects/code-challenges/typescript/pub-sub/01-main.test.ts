import { describe, it, assert, withSubmit } from "./unit_test";
import { type Move, type User, type Piece, doBattles, march } from "./main";

describe("doBattles", () => {
  const runCases = [
    {
      users: [
        {
          name: "Toussaint",
          pieces: [
            { location: "San Domingo", name: "Cavalry" },
            { location: "San Domingo", name: "Infantry" },
          ],
        },
        {
          name: "Napoleon",
          pieces: [
            { location: "France", name: "Infantry" },
            { location: "Russia", name: "Infantry" },
          ],
        },
        {
          name: "Washington",
          pieces: [{ location: "United States", name: "Artillery" }],
        },
      ],
      mv: {
        userName: "Toussaint",
        piece: { location: "United States", name: "Cavalry" },
      },
      expectedFightLocations: [
        { location: "United States", name: "Artillery" },
      ],
    },
  ];

  const submitCases = runCases.concat([
    {
      users: [
        {
          name: "Toussaint",
          pieces: [
            { location: "San Domingo", name: "Cavalry" },
            { location: "San Domingo", name: "Infantry" },
          ],
        },
        {
          name: "Napoleon",
          pieces: [
            { location: "France", name: "Infantry" },
            { location: "Russia", name: "Infantry" },
            { location: "United States", name: "Cavalry" },
          ],
        },
        {
          name: "Washington",
          pieces: [{ location: "United States", name: "Artillery" }],
        },
      ],
      mv: {
        userName: "Toussaint",
        piece: { location: "United States", name: "Cavalry" },
      },
      expectedFightLocations: [
        { location: "United States", name: "Cavalry" },
        { location: "United States", name: "Artillery" },
      ],
    },
  ]);

  const testCases = withSubmit ? submitCases : runCases;
  const skipped: number = submitCases.length - testCases.length;

  let passed: number = 0;
  let failed: number = 0;

  testCases.forEach((test, index) => {
    it(`Test ${index + 1}`, () => {
      const moves: Move[] = [];
      let mover: User | undefined;

      for (const u of test.users) {
        if (u.name === test.mv.userName) {
          mover = u;
          break;
        }
      }

      if (!mover) {
        console.error(
          `Test Failed: user with name ${test.mv.userName} not found`,
        );
        failed++;
        return;
      }

      march(mover, test.mv.piece, (move: Move) => moves.push(move));

      const output: Piece[] = doBattles(moves, test.users);
      try {
        assert.deepEqual(output, test.expectedFightLocations);
        console.log(`---------------------------------
          Test Passed:
            users:
          ${formatSlice(test.users)}
            move: ${JSON.stringify(test.mv)}
            =>
            expected battle pieces:
          ${formatSlice(test.expectedFightLocations)}
            actual battle pieces:
          ${formatSlice(output)}
          `);
        passed++;
      } catch {
        console.error(`---------------------------------
          Test Failed:
            users:
          ${formatSlice(test.users)}
            move: ${JSON.stringify(test.mv)}
            =>
            expected battle pieces:
          ${formatSlice(test.expectedFightLocations)}
            actual battle pieces:
          ${formatSlice(output)}
          `);
        failed++;
      }
    });
  });

  console.log("---------------------------------");
  if (skipped > 0) {
    console.log(`${passed} passed, ${failed} failed, ${skipped} skipped`);
  } else {
    console.log(`${passed} passed, ${failed} failed`);
  }
});

function formatSlice<T>(slice: T[]): string {
  if (!slice) return "null";
  return slice.map((item) => `* ${JSON.stringify(item)}`).join("\n");
}
