import * as z from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";

import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Textarea } from "@/components/ui/textarea";
import { Button } from "@/components/ui/button";
import { useCreateSecretMutation } from "@/feature/secrets/secrets.api";

const formScheme = z.object({
  message: z.string().min(3).max(800),
});

interface Props {
  onSubmit: (secretKey: string) => void;
}

export default function NewSecretForm({ onSubmit }: Props) {
  const form = useForm({
    resolver: zodResolver(formScheme),
    defaultValues: {
      message: "",
    },
  });

  const [createSecret] = useCreateSecretMutation();

  const handleSubmit = (data: z.infer<typeof formScheme>) => {
    createSecret({
      message: data.message,
    })
      .unwrap()
      .then((res) => {
        onSubmit(res.secretKey);
      });

    form.reset();
  };

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(handleSubmit)} className="space-y-8">
        <FormField
          control={form.control}
          name="message"
          render={({ field }) => (
            <FormItem>
              <FormLabel className="text-xl font-bold">Message</FormLabel>
              <FormControl>
                <Textarea
                  className="text-lg"
                  placeholder="..."
                  {...field}
                  autoComplete="off"
                  minLength={5}
                  maxLength={80}
                  rows={4}
                  autoFocus
                />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />
        <Button type="submit" size={"lg"}>
          <h6 className="text-lg">Create</h6>
        </Button>
      </form>
    </Form>
  );
}
