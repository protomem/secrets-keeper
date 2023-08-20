"use client";

import * as z from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";

import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";

const formScheme = z.object({
  message: z.string().min(3).max(800),
});

interface Props {
  onSubmit: (data: z.infer<typeof formScheme>) => void;
}

export default function NewSecret({ onSubmit: externalOnSubmit }: Props) {
  const form = useForm({
    resolver: zodResolver(formScheme),
    defaultValues: {
      message: "",
    },
  });

  const onSubmit = (data: z.infer<typeof formScheme>) => {
    console.log(data);
    localStorage.setItem("secret", data.message);
    externalOnSubmit(data);
    form.reset();
  };

  return (
    <Card className="w-[50rem] m-8">
      <CardHeader>
        <CardTitle className="text-center">New Secret</CardTitle>
      </CardHeader>
      <CardContent>
        <Form {...form}>
          <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
            <FormField
              control={form.control}
              name="message"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Message</FormLabel>
                  <FormControl>
                    <Textarea
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
            <Button type="submit">Create</Button>
          </form>
        </Form>
      </CardContent>
    </Card>
  );
}
