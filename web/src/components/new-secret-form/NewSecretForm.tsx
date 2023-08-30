import * as z from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { useCreateSecretMutation } from "@/feature/secrets/secrets.api";

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
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Input } from "@/components/ui/input";

const formScheme = z.object({
  message: z.string().min(3).max(800),
  ttl: z.number().min(0).max(3600),
  secretPhrase: z.string().min(3).max(80).optional(),
});

interface Props {
  onSubmit: (secretKey: string, withSecretPhrase: boolean) => void;
}

export default function NewSecretForm({ onSubmit }: Props) {
  const form = useForm({
    resolver: zodResolver(formScheme),
    defaultValues: {
      message: "",
      ttl: 0,
      secretPhrase: undefined,
    },
  });

  const [createSecret] = useCreateSecretMutation();

  const handleSubmit = (data: z.infer<typeof formScheme>) => {
    createSecret({
      message: data.message,
      ttl: data.ttl,
      secretPhrase: data.secretPhrase,
    })
      .unwrap()
      .then((res) => {
        onSubmit(res.secretKey, res.withSecretPhrase);
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
              <FormLabel className="text-xl">Message</FormLabel>
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

        <Accordion type="single" collapsible>
          <AccordionItem value="item-1" className="border-0">
            <AccordionTrigger className="text-xl">
              Advanced options
            </AccordionTrigger>
            <AccordionContent>
              <FormField
                control={form.control}
                name="ttl"
                render={({ field }) => (
                  <FormItem className="mx-2">
                    <FormLabel className="text-lg italic">TTL</FormLabel>
                    <Select
                      onValueChange={(value) => {
                        field.onChange(+value);
                      }}
                      defaultValue={field.value.toString()}
                    >
                      <FormControl>
                        <SelectTrigger>
                          <SelectValue />
                        </SelectTrigger>
                      </FormControl>

                      <SelectContent>
                        <SelectItem value="0">None</SelectItem>
                        <SelectItem value="1">1 hour</SelectItem>
                        <SelectItem value="3">3 hours</SelectItem>
                        <SelectItem value="6">6 hours</SelectItem>
                        <SelectItem value="12">12 hours</SelectItem>
                        <SelectItem value="24">1 day</SelectItem>
                        <SelectItem value="48">2 days</SelectItem>
                        <SelectItem value="120">5 days</SelectItem>
                        <SelectItem value="168">1 week</SelectItem>
                      </SelectContent>
                    </Select>

                    <FormMessage />
                  </FormItem>
                )}
              />

              <FormField
                control={form.control}
                name="secretPhrase"
                render={({ field }) => (
                  <FormItem className="mx-2 mt-4">
                    <FormLabel className="text-lg italic">
                      Secret phrase
                    </FormLabel>
                    <FormControl>
                      <Input placeholder="..." {...field} />
                    </FormControl>

                    <FormMessage />
                  </FormItem>
                )}
              />
            </AccordionContent>
          </AccordionItem>
        </Accordion>

        <Button type="submit" size={"lg"}>
          <h6 className="text-lg">Create</h6>
        </Button>
      </form>
    </Form>
  );
}
