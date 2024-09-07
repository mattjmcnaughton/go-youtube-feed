import { useState } from "react";

import { useForm } from "@mantine/form";
import {
  ActionIcon,
  Button,
  Box,
  CopyButton,
  Group,
  TextInput,
} from "@mantine/core";
import { IconCopy, IconCheck } from "@tabler/icons-react";

function FeedView() {
  const [feedURL, setFeedURL] = useState("");

  const form = useForm({
    mode: "uncontrolled",
    initialValues: {
      handle: "",
    },
  });

  async function formSubmit(values: { handle: string }) {
    setFeedURL(values.handle);
  }

  return (
    <>
      <Box>
        <h1>go-youtube-feed</h1>

        <form onSubmit={form.onSubmit(formSubmit)}>
          <TextInput
            label=""
            placeholder="handle"
            key={form.key("handle")}
            {...form.getInputProps("handle")}
          />

          <Button type="submit">Submit</Button>
        </form>

        {feedURL && (
          <Group>
            <p>{feedURL}</p>
            <CopyButton value={feedURL} timeout={2000}>
              {({ copied, copy }) => (
                <ActionIcon
                  color={copied ? "teal" : "gray"}
                  variant="subtle"
                  onClick={copy}
                >
                  {copied ? <IconCheck /> : <IconCopy />}
                </ActionIcon>
              )}
            </CopyButton>
          </Group>
        )}
      </Box>
    </>
  );
}

export default FeedView;
